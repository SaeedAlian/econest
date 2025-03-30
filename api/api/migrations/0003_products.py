from django.core.validators import MaxValueValidator, MinValueValidator
from django.db import migrations, models


class Migration(migrations.Migration):
    dependencies = [
        ("api", "0002_users"),
    ]

    operations = [
        migrations.CreateModel(
            name="Product",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "name",
                    models.CharField(
                        max_length=150,
                        unique=False,
                        blank=False,
                        null=False,
                    ),
                ),
                (
                    "slug",
                    models.CharField(
                        max_length=255,
                        unique=True,
                        blank=False,
                        null=False,
                    ),
                ),
                (
                    "price",
                    models.FloatField(
                        blank=False,
                        null=False,
                    ),
                ),
                (
                    "description",
                    models.CharField(
                        max_length=4095, unique=False, blank=True, null=False
                    ),
                ),
                ("created_at", models.DateTimeField(auto_now_add=True)),
                ("updated_at", models.DateTimeField(auto_now=True)),
            ],
            options={
                "db_table": "products",
            },
        ),
        migrations.CreateModel(
            name="ProductImage",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "image_name",
                    models.CharField(
                        max_length=255,
                        unique=True,
                        blank=False,
                        null=False,
                    ),
                ),
                ("is_main", models.BooleanField(default=False)),  # type: ignore
                (
                    "product",
                    models.ForeignKey(
                        to="Product",
                        on_delete=models.CASCADE,
                        db_column="product_id",
                    ),
                ),
            ],
            options={
                "db_table": "product_images",
            },
        ),
        migrations.CreateModel(
            name="ProductCategory",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "name",
                    models.CharField(
                        max_length=127,
                        unique=False,
                        blank=False,
                        null=False,
                    ),
                ),
                (
                    "parent_category",
                    models.ForeignKey(
                        to="ProductCategory",
                        on_delete=models.CASCADE,
                        db_column="product_category_id",
                    ),
                ),
                ("created_at", models.DateTimeField(auto_now_add=True)),
                ("updated_at", models.DateTimeField(auto_now=True)),
            ],
            options={
                "db_table": "product_categories",
            },
        ),
        migrations.AddField(
            model_name="Product",
            name="subcategory",
            field=models.ForeignKey(
                to="ProductCategory",
                on_delete=models.PROTECT,
                db_column="subcategory_id",
            ),
        ),
        migrations.CreateModel(
            name="ProductSpec",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "label",
                    models.CharField(
                        max_length=255,
                        unique=False,
                        blank=False,
                        null=False,
                    ),
                ),
                (
                    "value",
                    models.CharField(
                        max_length=255,
                        unique=False,
                        blank=False,
                        null=False,
                    ),
                ),
                (
                    "product",
                    models.ForeignKey(
                        to="Product",
                        on_delete=models.CASCADE,
                        db_column="product_id",
                    ),
                ),
            ],
            options={
                "db_table": "product_specs",
            },
        ),
        migrations.CreateModel(
            name="ProductTag",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "name",
                    models.CharField(
                        max_length=127, blank=False, null=False, unique=False
                    ),
                ),
                ("created_at", models.DateTimeField(auto_now_add=True)),
                ("updated_at", models.DateTimeField(auto_now=True)),
            ],
            options={
                "db_table": "product_tags",
            },
        ),
        migrations.CreateModel(
            name="ProductTagAssignment",
            fields=[
                (
                    "product",
                    models.ForeignKey(
                        to="Product",
                        on_delete=models.CASCADE,
                        db_column="product_id",
                        primary_key=True,
                        serialize=False,
                    ),
                ),
                (
                    "tag",
                    models.ForeignKey(
                        to="ProductTag",
                        on_delete=models.CASCADE,
                        db_column="tag_id",
                    ),
                ),
            ],
            options={
                "db_table": "product_tag_assignments",
                "unique_together": [("product", "tag")],
            },
        ),
        migrations.CreateModel(
            name="ProductAttribute",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "label",
                    models.CharField(
                        max_length=255,
                        unique=False,
                        blank=False,
                        null=False,
                    ),
                ),
                (
                    "product",
                    models.ForeignKey(
                        to="Product",
                        on_delete=models.CASCADE,
                        db_column="product_id",
                    ),
                ),
            ],
            options={
                "db_table": "product_attributes",
            },
        ),
        migrations.CreateModel(
            name="ProductAttributeOption",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "value",
                    models.CharField(
                        max_length=255,
                        unique=False,
                        blank=False,
                        null=False,
                    ),
                ),
                (
                    "attribute",
                    models.ForeignKey(
                        to="ProductAttribute",
                        on_delete=models.CASCADE,
                        db_column="attribute_id",
                    ),
                ),
            ],
            options={
                "db_table": "product_attribute_options",
            },
        ),
        migrations.CreateModel(
            name="ProductVariant",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "quantity",
                    models.IntegerField(
                        blank=False,
                        null=False,
                        default=0,  # type: ignore
                        validators=[
                            MinValueValidator(
                                0, message="quantity cannot be less than 0"
                            ),
                        ],
                    ),
                ),
                (
                    "product",
                    models.ForeignKey(
                        to="Product",
                        on_delete=models.CASCADE,
                        db_column="product_id",
                    ),
                ),
            ],
            options={
                "db_table": "product_variants",
            },
        ),
        migrations.CreateModel(
            name="ProductVariantOption",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "option",
                    models.ForeignKey(
                        to="ProductAttributeOption",
                        on_delete=models.CASCADE,
                        db_column="option_id",
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
                "db_table": "product_variant_options",
                "unique_together": [("variant", "option")],
            },
        ),
        migrations.CreateModel(
            name="ProductComment",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "scoring",
                    models.IntegerField(
                        blank=False,
                        null=False,
                        validators=[
                            MinValueValidator(1, message="score cannot be less than 1"),
                            MaxValueValidator(5, message="score cannot be more than 5"),
                        ],
                    ),
                ),
                (
                    "comment",
                    models.CharField(max_length=1023, blank=True, null=True),
                ),
                (
                    "product",
                    models.ForeignKey(
                        to="Product",
                        on_delete=models.CASCADE,
                        db_column="product_id",
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
                "db_table": "product_comments",
            },
        ),
    ]
