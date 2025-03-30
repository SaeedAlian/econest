from django.db import migrations, models


class Migration(migrations.Migration):
    dependencies = [
        ("api", "0004_permissions"),
    ]

    operations = [
        migrations.CreateModel(
            name="WalletTransaction",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                ("amount", models.FloatField(null=False, blank=False)),
                (
                    "tx_type",
                    models.CharField(max_length=20, null=False, blank=False),
                ),  # transaction_types enum
                (
                    "status",
                    models.CharField(
                        max_length=20, null=False, blank=False, default="pending"
                    ),
                ),  # transaction_status enum
                (
                    "wallet",
                    models.ForeignKey(
                        to="Wallet",
                        on_delete=models.CASCADE,
                        db_column="wallet_id",
                    ),
                ),
                ("created_at", models.DateTimeField(auto_now_add=True)),
            ],
            options={
                "db_table": "wallet_transactions",
            },
        ),
        migrations.RunSQL(
            """
            ALTER TABLE wallet_transactions 
            ALTER COLUMN tx_type TYPE transaction_types USING tx_type::transaction_types,
            ALTER COLUMN status TYPE transaction_status USING status::transaction_status;
            """
        ),
    ]
