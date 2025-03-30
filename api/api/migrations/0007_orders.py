from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ("api", "0006_vendors"),
    ]

    operations = [
        migrations.CreateModel(
            name="Order",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "total_price",
                    models.FloatField(
                        blank=False,
                        null=False,
                    ),
                ),
                (
                    "delivery_date",
                    models.DateField(
                        blank=False,
                        null=False,
                    ),
                ),
                (
                    "user",
                    models.ForeignKey(
                        to="User",
                        on_delete=models.DO_NOTHING,
                        db_column="user_id",
                    ),
                ),
                (
                    "transaction",
                    models.ForeignKey(
                        to="WalletTransaction",
                        unique=True,
                        on_delete=models.RESTRICT,
                        db_column="transaction_id",
                    ),
                ),
                ("verified", models.BooleanField(default=False)),  # type: ignore
                ("created_at", models.DateTimeField(auto_now_add=True)),
                ("updated_at", models.DateTimeField(auto_now=True)),
            ],
            options={
                "db_table": "orders",
            },
        ),
        migrations.CreateModel(
            name="OrderProductVariant",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "order",
                    models.ForeignKey(
                        to="Order",
                        on_delete=models.CASCADE,
                        db_column="order_id",
                    ),
                ),
                (
                    "variant",
                    models.ForeignKey(
                        to="ProductVariant",
                        on_delete=models.CASCADE,
                        db_column="variant_id",
                    ),
                ),
            ],
            options={
                "db_table": "order_product_variants",
            },
        ),
    ]
