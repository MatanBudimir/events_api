Requests example

What's needed?
- Postman

URLs and request body examples:
- POST: `http://localhost:8000/api/events/{event_id}/users`
    - User gets added to the specified event
    - Body: None
    - API gets the user from the JWT Token

- DELETE: `http://localhost:8000/api/meetings/{meeting_id}`
  - Meeting with the meeting_id gets deleted
  - Body: None

- GET: `http://localhost:8000/api/events/{event_id}/meetings`
  - Returns a list of all scheduled meetings for the given event
  - Body: None

- POST: `http://localhost:8000/api/meetings`
  - Creates a meeting with specified data
  - Example Body: 
  ```json
    {
    "name": "A test event",
    "description": "A test meeting description",
    "start_time": "2023-06-29T10:30:00Z",
    "end_time": "2023-06-29T12:00:21Z",
    "invitations": [
    {
    "user_id": 1
    },
    {
     "user_id": 2
    }
    ],
    "event_id": 1
    }
  ```
- GET: `http://localhost:8000/api/meetings/{meeting_id}/invitations`
  - Returns a list of all meeting invitations
- POST: `http://localhost:8000/api/meetings/{meeting_id}/invitations`
  - Creates an invitations to the meeting
  - Body: 
  ```json
    {
    "user_id": 1
    }
  ```
- DELETE: `http://localhost:8000/api/meetings/{meeting_id}/invitations/{invitation_id}`
  - Deletes the invite with the given invitation ID
- PUT: `http://localhost:8000/api/meetings/{meeting_id}/invitations/{invitation_id}`
  - Update the status of the invitations (accept/reject)
  - Body:
  ```json
        {
          "status": 1
        }
    ```
- POST: `http://localhost:8000/api/meetings/{meeting_id}/schedule`
  - Schedule the meeting with the given meeting ID
  - Meeting is scheduled only if everyone accepted the invite