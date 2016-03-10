# octav

builderscon web site and api server

# What's With The Name?

Eh, don't worry. Just a word that I came up with, and has absolutely no meaning.Feel free to suggest a better name

# Features

## MUST HAVE

* Administrators can register "organizers"
* Organizers can create/update/delete "conferences"
* Organizers can register "venues" (which has one or more "rooms")
* Organizers can accept/reject talk proposals
* Organizers can edit talk proposals, to set date/time, etc.
* Speakers can register their photo, bio. Must register email (not public, only used to send notices from organizers)
* Speakers can register themselves ("speakers")
* Speakers can submit talk proposals ("sessions"). Proposals can be either hidden or visible. If hidden, we need a "make schedule visible" button
* Conferences have "news feed", "Twitter/Social Network Feed display"
* Accepted talks show up in schedule
* Sessions can have, among other things, video urls and slide urls. These can be shown as a list. It can be grouped by tags, too.
* Sessions have Facebook/Twitter/etc buttons

## NICE TO HAVE

* Register users (attendees)
* Organizers can notify attendees via email or whatever else proper means
* Organizers can post to Facebook/Twitter/etc, via the official account.
* Attendees can send back feedback on sessions/conferences
* Attendees can vote on sessions to determine "best session"
* Sessions are announced 30 and 10 minutes before they are scheduled via Facebook/Twitter/etc
* Session details page show video and slides displayed in a synchronized fashion.

## MISCELLANEOUS

URL structure suggestions (just a thought, feel free to suggest better approaches).

Assume base url https://builderscon.io. "Main" builderscon site.

| name              | url pattern              | notes                                          |
|:------------------|:-------------------------|:-----------------------------------------------|
| main page         | /                        | latest conferences, links to videos, etc       |


Assume base url https://conf.builderscon.io. Conferences show up under this host

| name              | url pattern              | notes                                          |
|:------------------|:-------------------------|:-----------------------------------------------|
| conference page   | /tokyo                   | redirects to "latest" conference               |
| per-instance page | /tokyo/2017              | "2017" can be "2017-summer" or other subtitles |
| latest news       | /tokyo/2017/news         | |
| schedule/calendar | /tokyo/2017/schedule     | |
| session details   | /tokyo/2017/session/[id] | |
| speaker details   | /speaker/[id]            | use user.speaker\_bio |
| user details      | /user/[id]               | user can be an attendee, speaker, organizer, etc |


Admin site URL should be different, so let's assume base url https://admin.builderscon.io

| name                | url pattern         | notes                                          |
|:--------------------|:--------------------|:-----------------------------------------------|
| main page           | /dashboard          |                                                |
| register organizer  | /organizer/register | |
| register conference | /conference/create  | |
| TODO (Add more) | | |

# Development

## Useful Debugging Tips

### Enable Debug Prints

When running your tests, use the tag `debug0` (or debug). See [github.com/lestrrat/go-pdebug](https://github.com/lestrrat/go-pdebug) for details. |

## Useful Environment Variables

| Name | Description |
|:-----|:------------|
|OCTAV_TEST_DSN | DSN to use to connect to the database. See [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) for syntax details. |
|OCTAV_DB_TRACE | Enable to see SQL execution trace. See [github.com/shogo82148/go-sql-proxy](https://github.com/shogo82148/go-sql-proxy) for details. |
|PDEBUG_TRACE | Enable print debug mode. See [github.com/lestrrat/go-pdebug](https://github.com/lestrrat/go-pdebug) for details. |


# Running Tests

See "Useful Environemnt Variables", and "Useful Debugging Tips"

## 0. Drop the old database

The schema is still changing wildly. You probably want to flush it from
time to time

```
mysqladmin -uroot drop octav
```

## 1. Create a database

```
mysqladmin -uroot create octav
mysql -uroot octav < octav/sql/octav.sql
```

## 2. Run

```
cd octav
OCTAV_TEST_DSN='root:@:/octav?parseTime=true' go test .
```

Currently, the tests are still failing.
