package test

// Project SQL
const SelectAllProjectSQL = `SELECT * FROM "projects" WHERE "projects"."deleted_at" IS NULL`
const SelectOneByIdProjectSQL = `SELECT .* FROM "projects".*`
const InsertProjectSQL = `INSERT INTO "projects".*RETURNING "id","created_at","updated_at","deleted_at"`
const UpdateProjectSQL = `UPDATE "projects" SET`

// Filter Option SQL
const SelectAllFilterOptionSQL = `SELECT * FROM "project_filter_options" WHERE "project_filter_options"."deleted_at" IS NULL`
const SelectOneByIdFilterOptionSQL = `SELECT .* FROM "project_filter_options".*`
const InsertFilterOptionSQL = `INSERT INTO "project_filter_options".*RETURNING "id","created_at","updated_at","deleted_at"`
const UpdateFilterOptionSQL = `UPDATE "project_filter_options" SET`
