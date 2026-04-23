API: 

curl -X 'GET' \
  'http://localhost:8080/api/v1/board' \
  -H 'accept: application/json'


Response: 
  { "success": true,
  "data": {
    "columns": [
      {
        "id": "e468942b-0fb9-45fd-98d5-ad794fcde3f8",
        "name": "Backlog",
        "order": 0,
        "tasks": [
          {
            "id": "9316316b-87f4-47cc-af8f-dcfeb5b31968",
            "title": "TEST",
            "description": "XYZ",
            "column_id": "e468942b-0fb9-45fd-98d5-ad794fcde3f8",
            "assignee_id": "24c532b1-f7b3-414f-b6aa-2f84117b619d",
            "due_date": "2026-04-24T20:34:05+05:30",
            "created_by": "24c532b1-f7b3-414f-b6aa-2f84117b619d",
            "position": 0,
            "created_at": "2026-04-23T16:31:58+05:30",
            "updated_at": "2026-04-23T16:31:58+05:30"
          }
        ]
      },
      {
        "id": "46fcb0d1-ab10-4854-b1d8-85a9afe7d97b",
        "name": "To Do",
        "order": 1,
        "tasks": null
      },
      {
        "id": "dc5c550b-c551-43d6-8846-498e6145d494",
        "name": "In Progress",
        "order": 2,
        "tasks": null
      },
      {
        "id": "5e00b090-ee66-4a1c-9d89-2163354ef08c",
        "name": "Review",
        "order": 3,
        "tasks": null
      },
      {
        "id": "30419d9a-c30b-4bb8-ba1c-09fea5ad04d3",
        "name": "Testing",
        "order": 4,
        "tasks": null
      },
      {
        "id": "83e2e26a-3657-449a-9160-3ddfb00df15c",
        "name": "Done",
        "order": 5,
        "tasks": null
      },
      {
        "id": "00ce6ef4-d474-4823-8073-af9b4108632c",
        "name": "Blocked",
        "order": 6,
        "tasks": null
      },
      {
        "id": "bd3b7168-4452-408a-b805-0e70f1401ee0",
        "name": "Archived",
        "order": 7,
        "tasks": null
      }
    ]
  }
}

FRONTEND: http://localhost:3000/board

Uncaught TypeError: Cannot read properties of null (reading 'length')

If there is no task in the columns, then empty column should be shown with the title of the column.

Right now it is showing the above uncaught error.

Goal: To show all the columns along with the task and if there is no task empty column should be shown.

