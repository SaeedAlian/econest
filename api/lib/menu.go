package lib

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"

	"golang.org/x/term"

	db_manager "github.com/SaeedAlian/econest/api/db/manager"
	"github.com/SaeedAlian/econest/api/types"
)

const (
	CTRLC        = "\x03"
	KeyUp        = "\x1b[A"
	KeyDown      = "\x1b[B"
	KeyRight     = "\x1b[C"
	KeyLeft      = "\x1b[D"
	KeyBackspace = "\x08"
	KeyDelete    = "\x7f"
	KeyEnter     = "\r"
)

type OnClickHandler func(menu *Menu) error

type Menu struct {
	reader        *bufio.Reader
	pages         map[int]MenuPage
	currentPage   int
	currentOption int
	errorMessage  *string
	infoMessage   *string
	forceExit     bool

	User *types.User
}

type MenuPage struct {
	Label          string
	Options        []MenuPageOption
	PreviousMenuId *int
}

type MenuPageOption struct {
	Id      int
	Label   string
	OnClick OnClickHandler
}

func NewMenu(reader *bufio.Reader, db *db_manager.Manager, pages map[int]MenuPage) *Menu {
	return &Menu{
		reader:        reader,
		pages:         pages,
		currentPage:   1,
		currentOption: 1,
		forceExit:     false,
		User:          nil,
	}
}

func (m *Menu) Display() error {
	defer m.ShowCursor()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	m.ClearScreen()

	errorMessageCount := 0
	infoMessageCount := 0

	for {
		if m.forceExit {
			fmt.Println("Goodbye!")
			return nil
		}
		if m.errorMessage != nil {
			errorMessageCount++
		}

		if errorMessageCount > 5 {
			m.ClearErrorMessage()
			errorMessageCount = 0
		}

		if m.infoMessage != nil {
			infoMessageCount++
		}

		if infoMessageCount > 5 {
			m.ClearInfoMessage()
			infoMessageCount = 0
		}

		m.MoveCursorToTop()
		m.HideCursor()
		m.ClearScreen()

		page, err := m.GetCurrentPage()
		if err != nil {
			return err
		}

		options := page.Options
		if len(options) == 0 {
			return types.ErrPageHasNoOptions
		}

		var menuBuilder strings.Builder

		line := strings.Repeat("-", len(page.Label)+10)
		menuBuilder.WriteString(line + "\n")
		menuBuilder.WriteString(fmt.Sprintf("     %s     \n", page.Label))
		menuBuilder.WriteString(line + "\n\n\n")

		for _, o := range options {
			if o.Id == m.currentOption {
				menuBuilder.WriteString(fmt.Sprintf(">>> %d - %s\n", o.Id, o.Label))
			} else {
				menuBuilder.WriteString(fmt.Sprintf(" %d - %s\n", o.Id, o.Label))
			}
		}

		if m.errorMessage != nil {
			menuBuilder.WriteString(fmt.Sprintf("!!! Error: %s\n", *m.errorMessage))
		}

		if m.infoMessage != nil {
			menuBuilder.WriteString(fmt.Sprintf("--> %s\n", *m.infoMessage))
		}

		fmt.Print(menuBuilder.String())

		key, err := m.readKeyWithContext(ctx)
		if err != nil {
			if err == context.Canceled {
				fmt.Println("Goodbye!")
				return nil
			}
			return err
		}

		keepOn, err := m.onKeyPress(key)
		if err != nil {
			return err
		}
		if !keepOn {
			fmt.Println("Goodbye!")
			return nil
		}
	}
}

func (m *Menu) GetCurrentPage() (*MenuPage, error) {
	page, ok := m.pages[m.currentPage]
	if !ok {
		return nil, types.ErrPageNotFound
	}

	return &page, nil
}

func (m *Menu) ChangePage(pageId int) error {
	_, ok := m.pages[pageId]
	if !ok {
		return types.ErrPageNotFound
	}

	m.currentPage = pageId
	return nil
}

func (m *Menu) GetCurrentOption() (*MenuPageOption, error) {
	page, err := m.GetCurrentPage()
	if err != nil {
		return nil, err
	}

	opts := page.Options

	for _, o := range opts {
		if o.Id == m.currentOption {
			return &o, nil
		}
	}

	return nil, types.ErrPageOptionNotFound
}

func (m *Menu) IncreaseOption() {
	newOpt := m.currentOption + 1
	currPage, err := m.GetCurrentPage()
	if err != nil {
		return
	}

	if newOpt > len(currPage.Options) {
		newOpt = 1
	}

	m.currentOption = newOpt
}

func (m *Menu) DecreaseOption() {
	newOpt := m.currentOption - 1
	currPage, err := m.GetCurrentPage()
	if err != nil {
		return
	}

	if newOpt < 1 {
		newOpt = len(currPage.Options)
	}

	m.currentOption = newOpt
}

func (m *Menu) ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func (m *Menu) HideCursor() {
	fmt.Print("\033[?25l")
}

func (m *Menu) ShowCursor() {
	fmt.Print("\033[?25h")
}

func (m *Menu) MoveCursorToTop() {
	fmt.Print("\033[H")
}

func (m *Menu) ClearLine() {
	fmt.Print("\r\033[K")
}

func (m *Menu) SetErrorMessage(err error) {
	errStr := err.Error()
	m.errorMessage = &errStr
}

func (m *Menu) ClearErrorMessage() {
	m.errorMessage = nil
}

func (m *Menu) SetInfoMessage(message string) {
	m.infoMessage = &message
}

func (m *Menu) ClearInfoMessage() {
	m.infoMessage = nil
}

func (m *Menu) Exit() {
	m.forceExit = true
}

func (m *Menu) ReadKey() (string, error) {
	fd := int(os.Stdin.Fd())

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer term.Restore(fd, oldState)

	var buf [3]byte
	n, err := os.Stdin.Read(buf[:1])
	if err != nil || n == 0 {
		return "", err
	}

	if buf[0] == 0x1b {
		os.Stdin.Read(buf[1:3])
		return string(buf[:3]), nil
	}

	return string(buf[:1]), nil
}

func (m *Menu) PromptDateTime(message string) (*time.Time, error) {
	m.ShowCursor()
	fmt.Printf("%s (format YYYY-MM-DD,hh:mm:ss) ", message)
	inp, err := m.reader.ReadString('\n')
	inp = strings.TrimSpace(inp)
	if err != nil {
		return nil, err
	}

	splt := strings.Split(inp, ",")
	if len(splt) != 2 {
		return nil, types.ErrInvalidMenuInput
	}

	d := splt[0]
	t := splt[1]

	year, month, day, err := m.parseDateInput(d)
	if err != nil {
		return nil, err
	}
	hour, minute, second, err := m.parseTimeInput(t)
	if err != nil {
		return nil, err
	}

	newTime := time.Date(*year, time.Month(*month), *day, *hour, *minute, *second, 0, time.UTC)

	return &newTime, nil
}

func (m *Menu) PromptDate(message string) (*time.Time, error) {
	m.ShowCursor()
	fmt.Printf("%s (format YYYY-MM-DD) ", message)
	inp, err := m.reader.ReadString('\n')
	inp = strings.TrimSpace(inp)
	if err != nil {
		return nil, err
	}

	year, month, day, err := m.parseDateInput(inp)
	if err != nil {
		return nil, err
	}

	newTime := time.Date(*year, time.Month(*month), *day, 0, 0, 0, 0, time.UTC)

	return &newTime, nil
}

func (m *Menu) PromptBoolean(message string) (*bool, error) {
	m.ShowCursor()
	fmt.Print(message)
	inp, err := m.reader.ReadString('\n')
	inp = strings.TrimSpace(inp)
	if err != nil {
		return nil, err
	}

	parsed, err := strconv.Atoi(inp)
	if err != nil {
		return nil, err
	}

	var val bool

	switch parsed {
	case 0:
		{
			val = false
		}

	case 1:
		{
			val = true
		}

	default:
		{
			return nil, types.ErrInvalidMenuInput
		}
	}

	return &val, nil
}

func (m *Menu) PromptString(message string) (*string, error) {
	m.ShowCursor()
	fmt.Print(message)
	inp, err := m.reader.ReadString('\n')
	inp = strings.TrimSpace(inp)
	if err != nil {
		return nil, err
	}

	return &inp, nil
}

func (m *Menu) PromptPassword(message string) (*string, error) {
	m.ShowCursor()
	fmt.Print(message)
	inp, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println("")
	if err != nil {
		return nil, err
	}

	p := string(inp)
	p = strings.TrimSpace(p)

	return &p, nil
}

func (m *Menu) PromptFloat32(message string) (*float32, error) {
	m.ShowCursor()
	fmt.Print(message)
	inp, err := m.reader.ReadString('\n')
	inp = strings.TrimSpace(inp)
	if err != nil {
		return nil, err
	}

	parsed, err := strconv.ParseFloat(inp, 32)
	if err != nil {
		return nil, err
	}

	p := float32(parsed)

	return &p, nil
}

func (m *Menu) PromptFloat64(message string) (*float64, error) {
	m.ShowCursor()
	fmt.Print(message)
	inp, err := m.reader.ReadString('\n')
	inp = strings.TrimSpace(inp)
	if err != nil {
		return nil, err
	}

	parsed, err := strconv.ParseFloat(inp, 64)
	if err != nil {
		return nil, err
	}

	p := float64(parsed)

	return &p, nil
}

func (m *Menu) PromptInt(message string) (*int, error) {
	m.ShowCursor()
	fmt.Print(message)
	inp, err := m.reader.ReadString('\n')
	inp = strings.TrimSpace(inp)
	if err != nil {
		return nil, err
	}

	parsed, err := strconv.Atoi(inp)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func (m *Menu) PromptInt32(message string) (*int32, error) {
	m.ShowCursor()
	fmt.Print(message)
	inp, err := m.reader.ReadString('\n')
	inp = strings.TrimSpace(inp)
	if err != nil {
		return nil, err
	}

	parsed, err := strconv.ParseInt(inp, 10, 32)
	if err != nil {
		return nil, err
	}

	p := int32(parsed)

	return &p, nil
}

func (m *Menu) PromptInt64(message string) (*int64, error) {
	m.ShowCursor()
	fmt.Print(message)
	inp, err := m.reader.ReadString('\n')
	inp = strings.TrimSpace(inp)
	if err != nil {
		return nil, err
	}

	parsed, err := strconv.ParseInt(inp, 10, 64)
	if err != nil {
		return nil, err
	}

	p := int64(parsed)

	return &p, nil
}

func (m *Menu) parseDateInput(inp string) (year *int, month *int, day *int, err error) {
	splt := strings.Split(inp, "-")
	if len(splt) != 3 {
		return nil, nil, nil, types.ErrInvalidMenuInput
	}

	y := splt[0]
	mo := splt[1]
	d := splt[2]

	parsedYear, err := strconv.Atoi(y)
	if err != nil {
		return nil, nil, nil, types.ErrInvalidMenuInput
	}

	if parsedYear > 3000 || parsedYear < 1900 {
		return nil, nil, nil, types.ErrInvalidMenuInput
	}

	parsedMonth, err := strconv.Atoi(mo)
	if err != nil {
		return nil, nil, nil, types.ErrInvalidMenuInput
	}

	if parsedMonth > 12 || parsedMonth < 1 {
		return nil, nil, nil, types.ErrInvalidMenuInput
	}

	parsedDay, err := strconv.Atoi(d)
	if err != nil {
		return nil, nil, nil, types.ErrInvalidMenuInput
	}

	if parsedDay > 31 || parsedDay < 1 {
		return nil, nil, nil, types.ErrInvalidMenuInput
	}

	return &parsedYear, &parsedMonth, &parsedDay, nil
}

func (m *Menu) parseTimeInput(inp string) (hour *int, minute *int, second *int, err error) {
	splt := strings.Split(inp, ":")
	if len(splt) != 3 {
		return nil, nil, nil, types.ErrInvalidMenuInput
	}

	h := splt[0]
	min := splt[1]
	s := splt[2]

	parsedHour, err := strconv.Atoi(h)
	if err != nil {
		return nil, nil, nil, types.ErrInvalidMenuInput
	}

	if parsedHour > 23 || parsedHour < 0 {
		return nil, nil, nil, types.ErrInvalidMenuInput
	}

	parsedMinute, err := strconv.Atoi(min)
	if err != nil {
		return nil, nil, nil, types.ErrInvalidMenuInput
	}

	if parsedMinute > 59 || parsedMinute < 0 {
		return nil, nil, nil, types.ErrInvalidMenuInput
	}

	parsedSecond, err := strconv.Atoi(s)
	if err != nil {
		return nil, nil, nil, types.ErrInvalidMenuInput
	}

	if parsedSecond > 59 || parsedSecond < 0 {
		return nil, nil, nil, types.ErrInvalidMenuInput
	}

	return &parsedHour, &parsedMinute, &parsedSecond, nil
}

func (m *Menu) PromptStruct(target any, omitKeys []string) error {
	m.ShowCursor()
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return types.ErrTargetMustBeNonNilPtrToStruct
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return types.ErrTargetMustBePtrStruct
	}

	omitMap := make(map[string]bool)
	for _, k := range omitKeys {
		omitMap[k] = true
	}

	typ := val.Type()

	for i := range val.NumField() {
		field := typ.Field(i)
		if omitMap[field.Name] {
			continue
		}
		fieldVal := val.Field(i)
		if !fieldVal.CanSet() {
			continue
		}

		prompt := m.normalizeFieldName(field.Name)

		var (
			setVal reflect.Value
			err    error
		)

		switch fieldVal.Type().Kind() {
		case reflect.String:
			{
				var (
					p   *string
					err error
				)
				if strings.Contains(field.Name, "Password") {
					p, err = m.PromptPassword(prompt)
				} else {
					p, err = m.PromptString(prompt)
				}
				if err != nil {
					return err
				}
				setVal = reflect.ValueOf(*p)
			}

		case reflect.Bool:
			{
				p, err := m.PromptBoolean(prompt + " (0 = false, 1 = true) ")
				if err != nil {
					return err
				}
				setVal = reflect.ValueOf(*p)
			}

		case reflect.Int:
			{
				val, err := m.PromptInt(prompt)
				if err != nil {
					return err
				}

				setVal = reflect.ValueOf(val).Elem()
			}

		case reflect.Int32:
			{
				val, err := m.PromptInt32(prompt)
				if err != nil {
					return err
				}

				setVal = reflect.ValueOf(val).Elem()
			}

		case reflect.Int64:
			{
				val, err := m.PromptInt64(prompt)
				if err != nil {
					return err
				}

				setVal = reflect.ValueOf(val).Elem()
			}

		case reflect.Float32:
			{
				val, err := m.PromptFloat32(prompt)
				if err != nil {
					return err
				}
				setVal = reflect.ValueOf(val).Elem()
			}

		case reflect.Float64:
			{
				val, err := m.PromptFloat64(prompt)
				if err != nil {
					return err
				}
				setVal = reflect.ValueOf(val).Elem()
			}

		case reflect.Struct:
			{
				if fieldVal.Type() == reflect.TypeOf(time.Time{}) {
					p, err := m.PromptDate(prompt)
					if err != nil {
						return err
					}
					setVal = reflect.ValueOf(*p)
				} else {
					return types.ErrInvalidTargetFieldType
				}
			}

		case reflect.Slice:
			{
				arrKind := fieldVal.Type().Elem().Kind()
				sliceType := fieldVal.Type()
				elemType := sliceType.Elem()
				slice := reflect.MakeSlice(sliceType, 0, 0)

				ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
				defer stop()

				i := 1

				for {
					prmpt := fmt.Sprintf("%s (%d-th) (Ctrl+C to terminate) ", prompt, i)

					resultCh := make(chan reflect.Value)
					errCh := make(chan error)

					go func() {
						var (
							val any
							e   error
						)

						switch arrKind {
						case reflect.Int:
							val, e = m.PromptInt(prmpt)
						case reflect.Int32:
							val, e = m.PromptInt32(prmpt)
						case reflect.Int64:
							val, e = m.PromptInt64(prmpt)
						case reflect.Float32:
							val, e = m.PromptFloat32(prmpt)
						case reflect.Float64:
							val, e = m.PromptFloat64(prmpt)
						case reflect.String:
							val, e = m.PromptString(prmpt)
						case reflect.Bool:
							prmpt = prmpt + " (0 = false, 1 = true) "
							val, e = m.PromptBoolean(prmpt)
						case reflect.Struct:
							if elemType == reflect.TypeOf(time.Time{}) {
								val, e = m.PromptDate(prompt)
							} else {
								e = types.ErrInvalidTargetFieldType
							}
						}

						if e != nil {
							errCh <- e
							return
						}
						resultCh <- reflect.ValueOf(val).Elem()
					}()

					select {
					case <-ctx.Done():
						setVal = slice
						goto Done
					case err := <-errCh:
						return err
					case val := <-resultCh:
						slice = reflect.Append(slice, val)
						i++
					}
				}
			Done:
			}

		case reflect.Ptr:
			{
				elemKind := fieldVal.Type().Elem().Kind()
				var valPtr any
				switch elemKind {
				case reflect.String:
					if strings.Contains(field.Name, "Password") {
						valPtr, err = m.PromptPassword(prompt)
					} else {
						valPtr, err = m.PromptString(prompt)
					}
				case reflect.Bool:
					valPtr, err = m.PromptBoolean(prompt + " (0 = false, 1 = true) ")
				case reflect.Int:
					valPtr, err = m.PromptInt(prompt)
				case reflect.Int32:
					valPtr, err = m.PromptInt32(prompt)
				case reflect.Int64:
					valPtr, err = m.PromptInt64(prompt)
				case reflect.Float32:
					valPtr, err = m.PromptFloat32(prompt)
				case reflect.Float64:
					valPtr, err = m.PromptFloat64(prompt)
				case reflect.Struct:
					if fieldVal.Type().Elem() == reflect.TypeOf(time.Time{}) {
						valPtr, err = m.PromptDate(prompt)
					} else {
						return types.ErrInvalidTargetFieldType
					}
				}
				if err != nil {
					return err
				}
				if valPtr != nil {
					setVal = reflect.ValueOf(valPtr)
				}
			}

		default:
			{
				return types.ErrInvalidTargetFieldType
			}
		}

		if setVal.IsValid() {
			fieldVal.Set(setVal)
		}
	}
	return nil
}

func (m *Menu) normalizeFieldName(name string) string {
	var sb strings.Builder
	sb.WriteString("Enter ")
	for i, r := range name {
		if i > 0 && unicode.IsUpper(r) {
			sb.WriteRune(' ')
		}
		sb.WriteRune(r)
	}
	sb.WriteRune(':')
	str := sb.String()
	return str + " "
}

func (m *Menu) onKeyPress(key string) (bool, error) {
	switch key {
	case KeyBackspace, KeyDelete:
		{
			page, err := m.GetCurrentPage()
			if err != nil {
				return false, err
			}

			if page.PreviousMenuId != nil {
				err := m.ChangePage(*page.PreviousMenuId)
				if err != nil {
					return false, err
				}
			}

			return true, nil
		}
	case KeyRight:
		{
			return true, nil
		}
	case KeyLeft:
		{
			return true, nil
		}
	case KeyUp:
		{
			m.DecreaseOption()
			return true, nil
		}
	case KeyDown:
		{
			m.IncreaseOption()
			return true, nil
		}
	case KeyEnter:
		{
			opt, err := m.GetCurrentOption()
			if err != nil {
				return false, err
			}

			if opt.OnClick != nil {
				err := opt.OnClick(m)
				if err != nil {
					m.SetErrorMessage(err)
				}
			}

			return true, nil
		}

	case CTRLC:
		{
			return false, nil
		}
	}

	return true, nil
}

func (m *Menu) readKeyWithContext(ctx context.Context) (string, error) {
	keyCh := make(chan string)
	errCh := make(chan error)

	go func() {
		key, err := m.ReadKey()
		if err != nil {
			errCh <- err
			return
		}
		keyCh <- key
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-errCh:
		return "", err
	case key := <-keyCh:
		return key, nil
	}
}
