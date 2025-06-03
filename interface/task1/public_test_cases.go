package main

type TestCase struct {
	name     string
	input    string
	expected Message
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Парсинг простого JSON с разными типами данных",
		input: `{
	       "string": "text",
	       "number": 42,
	       "float": 3.14,
	       "bool": true,
	       "null": null
	   }`,
		expected: Message{
			"string": "text",
			"number": "42",
			"float":  "3.14",
			"bool":   "true",
			"null":   "",
		},
	},
	{
		name:  "Парсинг вложенного объекта",
		input: `{"user": {"name": "Alice", "age": 25}}`,
		expected: Message{
			"user": `{"name":"Alice","age":25}`,
		},
	},
	{
		name:  "Парсинг массива",
		input: `{"tags": ["one", "two", 3]}`,
		expected: Message{
			"tags": `["one","two",3]`,
		},
	},
	// Тесткейсы в помощь
	{
		name: "Парсинг комплексного JSON (пример из задания)",
		input: `{
	       "rawID":"d75641c7-958c-466c-8031-e484d9d85409",
	       "username": "admin",
	       "isAuthorized": true,
	       "tries": 4,
	       "session": {
	           "session_start": "2025-04-09T19:16:02+00:00",
	           "session_end": "2025-05-09T19:16:02+00:00",
	           "session_data": "U2VjdGlvbiAxLjEwLjMyIG9mICJkZSBGaW5pYnVzIEJvbm9ydW0gZXQgTWFsb3J1bSIsIHdyaXR0ZW4gYnkgQ2ljZXJvIGluIDQ1IEJDCiJTZWQgdXQgcGVyc3BpY2lhdGlzIHVuZGUgb21uaXMgaXN0ZSBuYXR1cyBlcnJvciBzaXQgdm9sdXB0YXRlbSBhY2N1c2FudGl1bSBkb2xvcmVtcXVlIGxhdWRhbnRpdW0sIHRvdGFtIHJlbSBhcGVyaWFtLCBlYXF1ZSBpcHNhIHF1YWUgYWIgaWxsbyBpbnZlbnRvcmUgdmVyaXRhdGlzIGV0IHF1YXNpIGFyY2hpdGVjdG8gYmVhdGFlIHZpdGFlIGRpY3RhIHN1bnQgZXhwbGljYWJvLiBOZW1vIGVuaW0gaXBzYW0gdm9sdXB0YXRlbSBxdWlhIHZvbHVwdGFzIHNpdCBhc3Blcm5hdHVyIGF1dCBvZGl0IGF1dCBmdWdpdCwgc2VkIHF1aWEgY29uc2VxdXVudHVyIG1hZ25pIGRvbG9yZXMgZW9zIHF1aSByYXRpb25lIHZvbHVwdGF0ZW0gc2VxdWkgbmVzY2l1bnQuIE5lcXVlIHBvcnJvIHF1aXNxdWFtIGVzdCwgcXVpIGRvbG9yZW0gaXBzdW0gcXVpYSBkb2xvciBzaXQgYW1ldCwgY29uc2VjdGV0dXIsIGFkaXBpc2NpIHZlbGl0LCBzZWQgcXVpYSBub24gbnVtcXVhbSBlaXVzIG1vZGkgdGVtcG9yYSBpbmNpZHVudCB1dCBsYWJvcmUgZXQgZG9sb3JlIG1hZ25hbSBhbGlxdWFtIHF1YWVyYXQgdm9sdXB0YXRlbS4gVXQgZW5pbSBhZCBtaW5pbWEgdmVuaWFtLCBxdWlzIG5vc3RydW0gZXhlcmNpdGF0aW9uZW0gdWxsYW0gY29ycG9yaXMgc3VzY2lwaXQgbGFib3Jpb3NhbSwgbmlzaSB1dCBhbGlxdWlkIGV4IGVhIGNvbW1vZGkgY29uc2VxdWF0dXI/IFF1aXMgYXV0ZW0gdmVsIGV1bSBpdXJlIHJlcHJlaGVuZGVyaXQgcXVpIGluIGVhIHZvbHVwdGF0ZSB2ZWxpdCBlc3NlIHF1YW0gbmloaWwgbW9sZXN0aWFlIGNvbnNlcXVhdHVyLCB2ZWwgaWxsdW0gcXVpIGRvbG9yZW0gZXVtIGZ1Z2lhdCBxdW8gdm9sdXB0YXMgbnVsbGEgcGFyaWF0dXI/Ig=="
	       },
	       "tags":["core","needResetPassword"]
	   }`,
		expected: Message{
			"rawID":        "d75641c7-958c-466c-8031-e484d9d85409",
			"username":     "admin",
			"isAuthorized": "true",
			"tries":        "4",
			"session":      `{"session_start":"2025-04-09T19:16:02+00:00","session_end":"2025-05-09T19:16:02+00:00","session_data":"U2VjdGlvbiAxLjEwLjMyIG9mICJkZSBGaW5pYnVzIEJvbm9ydW0gZXQgTWFsb3J1bSIsIHdyaXR0ZW4gYnkgQ2ljZXJvIGluIDQ1IEJDCiJTZWQgdXQgcGVyc3BpY2lhdGlzIHVuZGUgb21uaXMgaXN0ZSBuYXR1cyBlcnJvciBzaXQgdm9sdXB0YXRlbSBhY2N1c2FudGl1bSBkb2xvcmVtcXVlIGxhdWRhbnRpdW0sIHRvdGFtIHJlbSBhcGVyaWFtLCBlYXF1ZSBpcHNhIHF1YWUgYWIgaWxsbyBpbnZlbnRvcmUgdmVyaXRhdGlzIGV0IHF1YXNpIGFyY2hpdGVjdG8gYmVhdGFlIHZpdGFlIGRpY3RhIHN1bnQgZXhwbGljYWJvLiBOZW1vIGVuaW0gaXBzYW0gdm9sdXB0YXRlbSBxdWlhIHZvbHVwdGFzIHNpdCBhc3Blcm5hdHVyIGF1dCBvZGl0IGF1dCBmdWdpdCwgc2VkIHF1aWEgY29uc2VxdXVudHVyIG1hZ25pIGRvbG9yZXMgZW9zIHF1aSByYXRpb25lIHZvbHVwdGF0ZW0gc2VxdWkgbmVzY2l1bnQuIE5lcXVlIHBvcnJvIHF1aXNxdWFtIGVzdCwgcXVpIGRvbG9yZW0gaXBzdW0gcXVpYSBkb2xvciBzaXQgYW1ldCwgY29uc2VjdGV0dXIsIGFkaXBpc2NpIHZlbGl0LCBzZWQgcXVpYSBub24gbnVtcXVhbSBlaXVzIG1vZGkgdGVtcG9yYSBpbmNpZHVudCB1dCBsYWJvcmUgZXQgZG9sb3JlIG1hZ25hbSBhbGlxdWFtIHF1YWVyYXQgdm9sdXB0YXRlbS4gVXQgZW5pbSBhZCBtaW5pbWEgdmVuaWFtLCBxdWlzIG5vc3RydW0gZXhlcmNpdGF0aW9uZW0gdWxsYW0gY29ycG9yaXMgc3VzY2lwaXQgbGFib3Jpb3NhbSwgbmlzaSB1dCBhbGlxdWlkIGV4IGVhIGNvbW1vZGkgY29uc2VxdWF0dXI/IFF1aXMgYXV0ZW0gdmVsIGV1bSBpdXJlIHJlcHJlaGVuZGVyaXQgcXVpIGluIGVhIHZvbHVwdGF0ZSB2ZWxpdCBlc3NlIHF1YW0gbmloaWwgbW9sZXN0aWFlIGNvbnNlcXVhdHVyLCB2ZWwgaWxsdW0gcXVpIGRvbG9yZW0gZXVtIGZ1Z2lhdCBxdW8gdm9sdXB0YXMgbnVsbGEgcGFyaWF0dXI/Ig=="}`,
			"tags":         `["core","needResetPassword"]`,
		},
	},
}
