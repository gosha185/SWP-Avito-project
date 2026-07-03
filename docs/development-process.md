## Code Review Process

Every change to the project is developed in a separate feature branch and integrated into the protected `main` branch through a GitHub Pull Request. Direct pushes to `main` are disabled to ensure that all changes are reviewed and verified before becoming part of the production codebase.

The review process follows these steps:

1. A developer creates a branch from the corresponding GitHub Issue using the naming convention `<issue-number>-short-description`.
2. The implementation is completed together with any required tests and documentation updates.
3. A Pull Request is opened and linked to the related issue.
4. GitHub Actions automatically executes the CI pipeline.
5. Another team member reviews the Pull Request.
6. After the CI pipeline succeeds and at least one approval is received, the Pull Request is merged using the **Merge Commit** strategy.

The reviewer verifies several aspects of every change:

* implementation satisfies the Issue acceptance criteria;
* code follows the team's coding style and formatting conventions;
* architectural decisions remain consistent with the existing design;
* appropriate unit, integration, or quality requirement tests are included or updated;
* no secrets, credentials, or sensitive information are committed;
* documentation and `CHANGELOG.md` are updated when required.

The project follows the repository requirements established for the course:

* all changes, including documentation and configuration updates, are submitted through Pull Requests;
* Pull Request authors cannot approve their own changes;
* at least one approval from another team member is required before merging;
* Pull Requests are merged only after all required CI checks have passed successfully;
* Merge Commit is used as the merge strategy to preserve complete development history and traceability between Issues, Pull Requests, and releases.

This workflow improves code quality, knowledge sharing within the team, and reduces the risk of introducing regressions into the protected branch.

---

## Database and Migrations

The project uses **PostgreSQL** as its primary relational database. Database access is implemented using Go's standard `database/sql` package together with the `sqlx` library and the `lib/pq` PostgreSQL driver. The project intentionally uses raw SQL instead of an ORM to provide explicit control over database queries and schema evolution.

Database schema changes are managed using **golang-migrate**, which provides versioned, repeatable, and reversible migrations. Migration files are stored in the repository under the `migrations/` directory.

Each migration consists of two files:

* `<version>_<description>.up.sql` — applies schema changes;
* `<version>_<description>.down.sql` — reverts the corresponding changes.

This approach guarantees that every schema modification can be reproduced consistently across all development environments while also allowing safe rollback when necessary.

During local development, migrations are applied manually using the `golang-migrate` CLI:

```bash
migrate -path ./migrations -database "$DB_DSN" up
```

A rollback can be performed with:

```bash
migrate -path ./migrations -database "$DB_DSN" down 1
```

When the project is started with Docker Compose, SQL migration scripts are mounted into PostgreSQL's `/docker-entrypoint-initdb.d/` directory. This allows the database schema to be initialized automatically during the first creation of the database volume. Existing databases are not modified automatically, ensuring predictable schema management after the initial deployment.

The CI pipeline validates that migrations are syntactically correct and can be executed successfully, although migrations are not applied automatically to shared environments. Schema updates remain an explicit deployment step to avoid unintended database modifications.

The project does not use seed data. Instead, all business data is created through the application's API during normal operation or automated tests. This keeps the database state reproducible while avoiding discrepancies between predefined seed data and real application behavior.

Overall, the migration strategy provides:

* version-controlled database schema evolution;
* reproducible environment setup;
* safe rollback using paired migration files;
* clear separation between schema management and application logic;
* consistent database initialization across local development and containerized deployments.

