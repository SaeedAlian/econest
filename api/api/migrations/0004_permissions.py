from django.db import migrations, models


class Migration(migrations.Migration):
    dependencies = [
        ("api", "0003_products"),
    ]

    operations = [
        migrations.CreateModel(
            name="PermissionGroup",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "name",
                    models.CharField(
                        max_length=255, null=False, blank=False, unique=True
                    ),
                ),
                ("description", models.TextField(null=True, blank=True)),
                ("created_at", models.DateTimeField(auto_now_add=True)),
            ],
            options={
                "db_table": "permission_groups",
            },
        ),
        migrations.CreateModel(
            name="RolePermissionGroup",
            fields=[
                (
                    "role",
                    models.ForeignKey(
                        to="Role",
                        on_delete=models.CASCADE,
                        db_column="role_id",
                        primary_key=True,
                        serialize=False,
                    ),
                ),
                (
                    "permission_group",
                    models.ForeignKey(
                        to="PermissionGroup",
                        on_delete=models.CASCADE,
                        db_column="permission_group_id",
                    ),
                ),
            ],
            options={
                "db_table": "role_permission_groups",
                "unique_together": [("role", "permission_group")],
            },
        ),
        migrations.CreateModel(
            name="ResourcePermission",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "resource",
                    models.CharField(max_length=20, null=False, blank=False),
                ),  # resources enum
                ("description", models.TextField(null=True, blank=True)),
                ("created_at", models.DateTimeField(auto_now_add=True)),
                (
                    "group",
                    models.ForeignKey(
                        to="PermissionGroup",
                        on_delete=models.CASCADE,
                        db_column="group_id",
                    ),
                ),
            ],
            options={
                "db_table": "resource_permissions",
            },
        ),
        migrations.RunSQL(
            """
            ALTER TABLE resource_permissions 
            ALTER COLUMN resource TYPE resources USING resource::resources;
            """
        ),
        migrations.CreateModel(
            name="ActionPermission",
            fields=[
                ("id", models.AutoField(primary_key=True)),
                (
                    "action",
                    models.CharField(max_length=20, null=False, blank=False),
                ),  # actions enum
                ("description", models.TextField(null=True, blank=True)),
                ("created_at", models.DateTimeField(auto_now_add=True)),
                (
                    "group",
                    models.ForeignKey(
                        to="PermissionGroup",
                        on_delete=models.CASCADE,
                        db_column="group_id",
                    ),
                ),
            ],
            options={
                "db_table": "action_permissions",
            },
        ),
        migrations.RunSQL(
            """
            ALTER TABLE action_permissions 
            ALTER COLUMN action TYPE actions USING action::actions;
            """
        ),
    ]
