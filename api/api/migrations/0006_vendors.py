from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ("api", "0005_wallet_transactions"),
    ]

    operations = [
        migrations.CreateModel(
            name="Vendor",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "name",
                    models.CharField(
                        max_length=255,
                        unique=True,
                        blank=False,
                        null=False,
                    ),
                ),
                (
                    "description",
                    models.CharField(
                        max_length=1023,
                        unique=False,
                        blank=False,
                        null=False,
                    ),
                ),
                (
                    "owner",
                    models.ForeignKey(
                        to="User",
                        on_delete=models.CASCADE,
                        db_column="owner_id",
                    ),
                ),
                ("verified", models.BooleanField(default=False)),  # type: ignore
                ("created_at", models.DateTimeField(auto_now_add=True)),
                ("updated_at", models.DateTimeField(auto_now=True)),
            ],
            options={
                "db_table": "vendors",
            },
        ),
        migrations.AlterField(
            model_name="Address",
            name="user",
            field=models.ForeignKey(
                to="User",
                on_delete=models.CASCADE,
                related_name="addresses",
                db_column="user_id",
                null=True,
                blank=True,
            ),
        ),
        migrations.AddField(
            model_name="Address",
            name="vendor",
            field=models.ForeignKey(
                to="Vendor",
                on_delete=models.CASCADE,
                related_name="addresses",
                db_column="vendor_id",
                null=True,
                blank=True,
            ),
        ),
        migrations.AddConstraint(
            model_name="Address",
            constraint=models.CheckConstraint(
                check=models.Q(user__isnull=False) | models.Q(vendor__isnull=False),
                name="address_has_owner",
            ),
        ),
        migrations.AddConstraint(
            model_name="Address",
            constraint=models.CheckConstraint(
                check=~models.Q(user__isnull=False, vendor__isnull=False),
                name="address_single_owner",
            ),
        ),
        migrations.AlterField(
            model_name="PhoneNumber",
            name="user",
            field=models.ForeignKey(
                to="User",
                on_delete=models.CASCADE,
                related_name="phone_numbers",
                db_column="user_id",
                null=True,
                blank=True,
            ),
        ),
        migrations.AddField(
            model_name="PhoneNumber",
            name="vendor",
            field=models.ForeignKey(
                to="Vendor",
                on_delete=models.CASCADE,
                related_name="phone_numbers",
                db_column="vendor_id",
                null=True,
                blank=True,
            ),
        ),
        migrations.AddConstraint(
            model_name="PhoneNumber",
            constraint=models.CheckConstraint(
                check=models.Q(user__isnull=False) | models.Q(vendor__isnull=False),
                name="phone_number_has_owner",
            ),
        ),
        migrations.AddConstraint(
            model_name="PhoneNumber",
            constraint=models.CheckConstraint(
                check=~models.Q(user__isnull=False, vendor__isnull=False),
                name="phone_number_single_owner",
            ),
        ),
        migrations.CreateModel(
            name="VendorProduct",
            fields=[
                (
                    "vendor",
                    models.ForeignKey(
                        to="Vendor",
                        on_delete=models.RESTRICT,
                        db_column="vendor_id",
                        primary_key=True,
                        serialize=False,
                    ),
                ),
                (
                    "product",
                    models.ForeignKey(
                        to="Product",
                        on_delete=models.RESTRICT,
                        db_column="product_id",
                    ),
                ),
            ],
            options={
                "db_table": "vendor_products",
                "unique_together": [("product", "vendor")],
            },
        ),
    ]
