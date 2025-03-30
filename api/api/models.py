from django.db import models
from django.core.validators import MinValueValidator, RegexValidator, MaxValueValidator
from django.core.exceptions import ValidationError


class Role(models.Model):
    id = models.AutoField(primary_key=True)
    name = models.CharField(max_length=255, unique=True, blank=False, null=False)
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        db_table = "roles"


class User(models.Model):
    id = models.AutoField(primary_key=True)
    username = models.CharField(
        max_length=150,
        unique=True,
        blank=False,
        null=False,
        validators=[
            RegexValidator(
                regex=r"^[a-zA-Z0-9]+[a-zA-Z0-9_-]*",
                message="Invalid username. It must start with a letter or number and can only contain letters, numbers, underscores, and hyphens.",
                code=400,
            )
        ],
    )
    email = models.CharField(
        max_length=255,
        unique=True,
        blank=False,
        null=False,
        validators=[
            RegexValidator(
                regex=r"^(?i)([a-z0-9._%+-]+)@(gmail\.com|googlemail\.com|yahoo\.com|yahoo\.co\.uk|protonmail\.com|proton\.me|outlook\.com|hotmail\.com|icloud\.com|mail\.com|aol\.com|zoho\.com)$",
                message="Invalid email address. Please use a valid email from supported domains: Gmail, Googlemail, Yahoo, ProtonMail, Outlook, Hotmail, iCloud, Mail.com, AOL, or Zoho.",
                code=400,
            )
        ],
    )
    email_verified = models.BooleanField(default=False)  # type: ignore
    password = models.CharField(max_length=256, unique=False, blank=False, null=False)
    full_name = models.CharField(max_length=255, unique=False, blank=True, null=True)
    birth_date = models.DateField()
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    role = models.ForeignKey(
        Role,
        on_delete=models.PROTECT,
        db_column="role_id",
    )

    class Meta:
        db_table = "users"


class Wallet(models.Model):
    id = models.AutoField(primary_key=True)
    user = models.ForeignKey(User, on_delete=models.CASCADE, db_column="user_id")
    balance = models.FloatField(
        default=0,  # type: ignore
        null=False,
        blank=False,
        validators=[
            MinValueValidator(0, message="Wallet balance cannot be less than 0")
        ],
    )
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        db_table = "wallets"


class WalletTransaction(models.Model):
    class TransactionType(models.TextChoices):
        DEPOSIT = "deposit", "Deposit"
        WITHDRAW = "withdraw", "Withdraw"
        PURCHASE = "purchase", "Purchase"
        SALE = "sale", "Sale"

    class Status(models.TextChoices):
        PENDING = "pending", "Pending"
        SUCCESSFUL = "successful", "Successful"
        FAILED = "failed", "Failed"

    id = models.AutoField(primary_key=True)
    wallet = models.ForeignKey(Wallet, on_delete=models.CASCADE, db_column="wallet_id")
    tx_type = models.CharField(
        max_length=20, null=False, blank=False, choices=TransactionType.choices
    )
    amount = models.FloatField(
        null=False,
        blank=False,
        validators=[
            MinValueValidator(0, message="Transaction amount cannot be less than 0")
        ],
    )
    status = models.CharField(
        max_length=20,
        null=False,
        blank=False,
        choices=Status.choices,
        default=Status.PENDING,
    )
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        db_table = "wallet_transactions"


class PermissionGroup(models.Model):
    id = models.AutoField(primary_key=True)
    name = models.CharField(max_length=255, null=False, blank=False, unique=True)
    description = models.TextField(null=True, blank=True)
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        db_table = "permission_groups"


class RolePermissionGroup(models.Model):
    role = models.ForeignKey(
        Role,
        on_delete=models.CASCADE,
        db_column="role_id",
    )
    permission_group = models.ForeignKey(
        PermissionGroup, on_delete=models.CASCADE, db_column="permission_group_id"
    )

    class Meta:
        db_table = "role_permission_groups"
        unique_together = ("role", "permission_group")


class ResourcePermission(models.Model):
    class Resource(models.TextChoices):
        PRODUCTS = "products", "Products"
        ORDERS = "orders", "Orders"
        USERS = "users", "Users"
        PERMISSIONS = "permissions", "Permissions"

    id = models.AutoField(primary_key=True)
    resource = models.CharField(
        max_length=63, null=False, blank=False, choices=Resource.choices
    )
    description = models.TextField(null=True, blank=True)
    group = models.ForeignKey(
        PermissionGroup, on_delete=models.CASCADE, db_column="group_id"
    )
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        db_table = "resource_permissions"


class ActionPermission(models.Model):
    class Action(models.TextChoices):
        CAN_ADD_PRODUCT = "can_add_product", "Can Add Product"
        CAN_UPDATE_PRODUCT = "can_update_product", "Can Update Product"
        CAN_DELETE_PRODUCT = "can_delete_product", "Can Delete Product"
        CAN_ADD_VENDOR = "can_add_vendor", "Can Add Vendor"
        CAN_UPDATE_VENDOR = "can_update_vendor", "Can Update Vendor"
        CAN_DELETE_VENDOR = "can_delete_vendor", "Can Delete Vendor"
        CAN_BAN_USER = "can_ban_user", "Can Ban User"
        CAN_UNBAN_USER = "can_unban_user", "Can Unban User"
        CAN_ADD_PRODUCT_TAG = "can_add_product_tag", "Can Add Product Tag"
        CAN_DELETE_PRODUCT_TAG = "can_delete_product_tag", "Can Delete Product Tag"
        CAN_ADD_PRODUCT_CATEGORY = (
            "can_add_product_category",
            "Can Add Product Category",
        )
        CAN_DELETE_PRODUCT_CATEGORY = (
            "can_delete_product_category",
            "Can Delete Product Category",
        )
        CAN_DELETE_PRODUCT_COMMENT = (
            "can_delete_product_comment",
            "Can Delete Product Comment",
        )
        CAN_ADD_ROLE = "can_add_role", "Can Add Role"
        CAN_DELETE_ROLE = "can_delete_role", "Can Delete Role"
        CAN_MODIFY_ROLE = "can_modify_role", "Can Modify Role"
        CAN_ADD_PERMISSION_GROUP = (
            "can_add_permission_group",
            "Can Add Permission Group",
        )
        CAN_DELETE_PERMISSION_GROUP = (
            "can_delete_permission_group",
            "Can Delete Permission Group",
        )
        CAN_MODIFY_PERMISSION_GROUP = (
            "can_modify_permission_group",
            "Can Modify Permission Group",
        )

    id = models.AutoField(primary_key=True)
    action = models.CharField(
        max_length=63, null=False, blank=False, choices=Action.choices
    )
    description = models.TextField(null=True, blank=True)
    group = models.ForeignKey(
        PermissionGroup, on_delete=models.CASCADE, db_column="group_id"
    )
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        db_table = "action_permissions"


class ProductCategory(models.Model):
    id = models.AutoField(primary_key=True)
    name = models.CharField(
        max_length=127,
        unique=False,
        blank=False,
        null=False,
    )
    parent_category = models.ForeignKey(
        "self",
        on_delete=models.CASCADE,
        db_column="product_category_id",
        null=True,
        blank=True,
    )
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        db_table = "product_categories"


class Product(models.Model):
    id = models.AutoField(primary_key=True)
    name = models.CharField(
        max_length=150,
        unique=False,
        blank=False,
        null=False,
    )
    slug = models.CharField(
        max_length=255,
        unique=True,
        blank=False,
        null=False,
    )
    price = models.FloatField(
        blank=False,
        null=False,
    )
    description = models.CharField(
        max_length=4095, unique=False, blank=True, null=False
    )
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    subcategory = models.ForeignKey(
        ProductCategory, on_delete=models.PROTECT, db_column="subcategory_id"
    )

    class Meta:
        db_table = "products"


class ProductImage(models.Model):
    id = models.AutoField(primary_key=True)
    image_name = models.CharField(
        max_length=255,
        unique=True,
        blank=False,
        null=False,
    )
    is_main = models.BooleanField(default=False)  # type: ignore
    product = models.ForeignKey(
        Product, on_delete=models.CASCADE, db_column="product_id"
    )

    class Meta:
        db_table = "product_images"


class ProductSpec(models.Model):
    id = models.AutoField(primary_key=True)
    label = models.CharField(
        max_length=255,
        unique=False,
        blank=False,
        null=False,
    )
    value = models.CharField(
        max_length=255,
        unique=False,
        blank=False,
        null=False,
    )
    product = models.ForeignKey(
        Product, on_delete=models.CASCADE, db_column="product_id"
    )

    class Meta:
        db_table = "product_specs"


class ProductTag(models.Model):
    id = models.AutoField(primary_key=True)
    name = models.CharField(max_length=127, blank=False, null=False, unique=False)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        db_table = "product_tags"


class ProductTagAssignment(models.Model):
    product = models.ForeignKey(
        Product,
        on_delete=models.CASCADE,
        db_column="product_id",
    )
    tag = models.ForeignKey(ProductTag, on_delete=models.CASCADE, db_column="tag_id")

    class Meta:
        db_table = "product_tag_assignments"
        unique_together = ("product", "tag")


class ProductAttribute(models.Model):
    id = models.AutoField(primary_key=True)
    label = models.CharField(
        max_length=255,
        unique=False,
        blank=False,
        null=False,
    )
    product = models.ForeignKey(
        Product, on_delete=models.CASCADE, db_column="product_id"
    )

    class Meta:
        db_table = "product_attributes"


class ProductAttributeOption(models.Model):
    id = models.AutoField(primary_key=True)
    value = models.CharField(
        max_length=255,
        unique=False,
        blank=False,
        null=False,
    )
    attribute = models.ForeignKey(
        ProductAttribute, on_delete=models.CASCADE, db_column="attribute_id"
    )

    class Meta:
        db_table = "product_attribute_options"


class ProductVariant(models.Model):
    id = models.AutoField(primary_key=True)
    quantity = models.IntegerField(
        blank=False,
        null=False,
        default=0,  # type: ignore
        validators=[
            MinValueValidator(0, message="quantity cannot be less than 0"),
        ],
    )
    product = models.ForeignKey(
        Product, on_delete=models.CASCADE, db_column="product_id"
    )

    class Meta:
        db_table = "product_variants"


class ProductVariantOption(models.Model):
    id = models.AutoField(primary_key=True)
    variant = models.ForeignKey(
        ProductVariant, on_delete=models.CASCADE, db_column="variant_id"
    )
    option = models.ForeignKey(
        ProductAttributeOption, on_delete=models.CASCADE, db_column="option_id"
    )

    class Meta:
        db_table = "product_variant_options"
        unique_together = ("variant", "option")


class ProductComment(models.Model):
    id = models.AutoField(primary_key=True)
    scoring = models.IntegerField(
        blank=False,
        null=False,
        validators=[
            MinValueValidator(1, message="score cannot be less than 1"),
            MaxValueValidator(5, message="score cannot be more than 5"),
        ],
    )
    comment = models.CharField(max_length=1023, blank=True, null=True)
    product = models.ForeignKey(
        Product, on_delete=models.CASCADE, db_column="product_id"
    )
    user = models.ForeignKey(User, on_delete=models.CASCADE, db_column="user_id")
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        db_table = "product_comments"


class Vendor(models.Model):
    id = models.AutoField(primary_key=True)
    name = models.CharField(
        max_length=255,
        unique=True,
        blank=False,
        null=False,
    )
    description = models.CharField(
        max_length=1023,
        unique=False,
        blank=False,
        null=False,
    )
    owner = models.ForeignKey(User, on_delete=models.CASCADE, db_column="owner_id")
    verified = models.BooleanField(default=False)  # type: ignore
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        db_table = "vendors"


class Address(models.Model):
    id = models.AutoField(primary_key=True)
    state = models.CharField(
        max_length=127,
        unique=False,
        blank=False,
        null=False,
    )
    city = models.CharField(
        max_length=127,
        unique=False,
        blank=False,
        null=False,
    )
    street = models.CharField(
        max_length=255,
        unique=False,
        blank=False,
        null=False,
    )
    zipcode = models.CharField(
        max_length=127,
        unique=False,
        blank=False,
        null=False,
        validators=[
            RegexValidator(
                regex=r"^[0-9]+$",
                message="Invalid zipcode. It must only contain digits.",
                code=400,
            )
        ],
    )
    details = models.TextField(
        max_length=1023,
        unique=False,
        blank=True,
        null=True,
    )
    user = models.ForeignKey(
        User,
        on_delete=models.CASCADE,
        related_name="addresses",
        db_column="user_id",
        null=True,
        blank=True,
    )
    vendor = models.ForeignKey(
        Vendor,
        on_delete=models.CASCADE,
        related_name="addresses",
        db_column="vendor_id",
        null=True,
        blank=True,
    )
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        db_table = "addresses"
        constraints = [
            models.CheckConstraint(
                check=models.Q(user__isnull=False) | models.Q(vendor__isnull=False),
                name="address_has_owner",
            ),
            models.CheckConstraint(
                check=~models.Q(user__isnull=False, vendor__isnull=False),
                name="address_single_owner",
            ),
        ]

    def clean(self):
        if not self.user and not self.vendor:
            raise ValidationError("Address must belong to either a user or vendor")
        if self.user and self.vendor:
            raise ValidationError("Address cannot belong to both user and vendor")


class PhoneNumber(models.Model):
    id = models.AutoField(primary_key=True)
    country_code = models.CharField(
        max_length=5,
        unique=False,
        blank=False,
        null=False,
        validators=[
            RegexValidator(
                regex=r"\+[0-9]{1,4}",
                message="Invalid country code. It must start with '+' followed by 1 to 4 digits.",
                code=400,
            )
        ],
    )
    number = models.CharField(
        max_length=15,
        unique=True,
        blank=False,
        null=False,
        validators=[
            RegexValidator(
                regex=r"^[0-9]+$",
                message="Invalid phone number. It must only contain digits.",
                code=400,
            )
        ],
    )
    user = models.ForeignKey(
        User,
        on_delete=models.CASCADE,
        related_name="phone_numbers",
        db_column="user_id",
        null=True,
        blank=True,
    )
    vendor = models.ForeignKey(
        "Vendor",
        on_delete=models.CASCADE,
        related_name="phone_numbers",
        db_column="vendor_id",
        null=True,
        blank=True,
    )
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        db_table = "phonenumbers"
        constraints = [
            models.CheckConstraint(
                check=models.Q(user__isnull=False) | models.Q(vendor__isnull=False),
                name="phone_number_has_owner",
            ),
            models.CheckConstraint(
                check=~models.Q(user__isnull=False, vendor__isnull=False),
                name="phone_number_single_owner",
            ),
        ]

    def clean(self):
        if not self.user and not self.vendor:
            raise ValidationError("Phone number must belong to either a user or vendor")
        if self.user and self.vendor:
            raise ValidationError("Phone number cannot belong to both user and vendor")


class VendorProduct(models.Model):
    vendor = models.ForeignKey(
        "Vendor",
        on_delete=models.RESTRICT,
        db_column="vendor_id",
    )
    product = models.ForeignKey(
        "Product", on_delete=models.RESTRICT, db_column="product_id"
    )

    class Meta:
        db_table = "vendor_products"
        unique_together = ("product", "vendor")


class Order(models.Model):
    id = models.AutoField(primary_key=True)
    total_price = models.FloatField(
        blank=False,
        null=False,
    )
    delivery_date = models.DateField(
        blank=False,
        null=False,
    )
    user = models.ForeignKey(User, on_delete=models.DO_NOTHING, db_column="user_id")
    transaction = models.OneToOneField(
        "WalletTransaction",
        on_delete=models.RESTRICT,
        db_column="transaction_id",
        unique=True,
    )
    verified = models.BooleanField(default=False)  # type: ignore
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        db_table = "orders"


class OrderProductVariant(models.Model):
    id = models.AutoField(primary_key=True)
    order = models.ForeignKey("Order", on_delete=models.CASCADE, db_column="order_id")
    variant = models.ForeignKey(
        "ProductVariant", on_delete=models.CASCADE, db_column="variant_id"
    )

    class Meta:
        db_table = "order_product_variants"
