from django.core.validators import RegexValidator, MinValueValidator
from django.db import migrations, models


class Migration(migrations.Migration):
    dependencies = [
        ("api", "0001_initial"),
    ]

    operations = [
        migrations.CreateModel(
            name="User",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "username",
                    models.CharField(
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
                    ),
                ),
                (
                    "email",
                    models.CharField(
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
                    ),
                ),
                ("email_verified", models.BooleanField(default=False)),  # type: ignore
                (
                    "password",
                    models.CharField(
                        max_length=256, unique=False, blank=False, null=False
                    ),
                ),
                (
                    "full_name",
                    models.CharField(
                        max_length=255, unique=False, blank=True, null=True
                    ),
                ),
                ("birth_date", models.DateTimeField()),
                ("created_at", models.DateTimeField(auto_now_add=True)),
                ("updated_at", models.DateTimeField(auto_now=True)),
            ],
            options={
                "db_table": "users",
            },
        ),
        migrations.CreateModel(
            name="PhoneNumber",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "country_code",
                    models.CharField(
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
                    ),
                ),
                (
                    "number",
                    models.CharField(
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
                    ),
                ),
                (
                    "user",
                    models.ForeignKey(
                        to="User",
                        on_delete=models.CASCADE,
                        related_name="phone_numbers",
                        db_column="user_id",
                    ),
                ),
                ("verified", models.BooleanField(default=False)),  # type: ignore
                ("created_at", models.DateTimeField(auto_now_add=True)),
                ("updated_at", models.DateTimeField(auto_now=True)),
            ],
            options={
                "db_table": "phonenumbers",
            },
        ),
        migrations.CreateModel(
            name="Address",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "state",
                    models.CharField(
                        max_length=127,
                        unique=False,
                        blank=False,
                        null=False,
                    ),
                ),
                (
                    "city",
                    models.CharField(
                        max_length=127,
                        unique=False,
                        blank=False,
                        null=False,
                    ),
                ),
                (
                    "street",
                    models.CharField(
                        max_length=255,
                        unique=False,
                        blank=False,
                        null=False,
                    ),
                ),
                (
                    "zipcode",
                    models.CharField(
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
                    ),
                ),
                (
                    "details",
                    models.CharField(
                        max_length=1023,
                        unique=False,
                        blank=True,
                        null=True,
                    ),
                ),
                (
                    "user",
                    models.ForeignKey(
                        to="User",
                        on_delete=models.CASCADE,
                        related_name="addresses",
                        db_column="user_id",
                    ),
                ),
                ("created_at", models.DateTimeField(auto_now_add=True)),
                ("updated_at", models.DateTimeField(auto_now=True)),
            ],
            options={
                "db_table": "addresses",
            },
        ),
        migrations.CreateModel(
            name="Role",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "name",
                    models.CharField(
                        max_length=255, unique=True, blank=False, null=False
                    ),
                ),
                ("created_at", models.DateTimeField(auto_now_add=True)),
            ],
            options={
                "db_table": "roles",
            },
        ),
        migrations.AddField(
            model_name="User",
            name="role",
            field=models.ForeignKey(
                to="Role",
                on_delete=models.PROTECT,
                db_column="role_id",
            ),
        ),
        migrations.CreateModel(
            name="Wallet",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "balance",
                    models.FloatField(
                        default=0,  # type: ignore
                        null=False,
                        blank=False,
                        validators=[
                            MinValueValidator(
                                0, message="Wallet balance cannot be less than 0"
                            )
                        ],
                    ),
                ),
                (
                    "user",
                    models.ForeignKey(
                        to="User",
                        on_delete=models.CASCADE,
                        db_column="user_id",
                    ),
                ),
                ("created_at", models.DateTimeField(auto_now_add=True)),
                ("updated_at", models.DateTimeField(auto_now=True)),
            ],
            options={
                "db_table": "wallets",
            },
        ),
    ]
