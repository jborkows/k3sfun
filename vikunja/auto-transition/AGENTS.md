Use https://vikunja.borkowskij.com/api/v1/docs.json for openapi api documentation
Compile and run in context of current directory.
Use queries.go to put read only commands.
See that BucketMapping in definitions.go has mapping of name to id.
Keep only meaningful comments, prefer creating functions to comments.

Files:
- commands.go modification api calls
- configuration.go
- definitions.go dtos
- queries.go read only api calls
- task.go contain main logic loop inside of function Run


The Kanban view requires the view-specific bucket endpoint to properly move tasks between buckets.
Environment variables with api token are available.
