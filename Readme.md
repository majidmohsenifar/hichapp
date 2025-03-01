# HICH-APP test project
this project is a simple system to create polls and vote to them.


### How to run the project
you need to have docker and docker-compose installed on your machine.
clone the project

run ```docker-compose up -d```

open the swagger in address ```http://localhost:8000/swagger/index.html```

you can access to prometheus UI using address ```http://localhost:9090```


### project structure
- the project has 3 layer consist of handlers(api), service, and repo.
service is the place that we store the bussiness logic.


### test
-integration tests are located in tests directory
- integration tests cover most of the scenarios for all the APIs, so the unit tests are not implemented.

to run test "go test ./test/..."


### database structure
- polls, options, tags are different tables.
- the relation between poll and tag is many to many with pivot table poll_tags.
- the relation between poll and options is one to many. 
- there is a table named votes which store the votes and skip of the polls by user.
we can consider the votes table as interaction of users with polls whether vote or skip.
- the count of votes for each options is stored in the options table, this can help us to not
do heavy query for counting the votes for stats api.


### technical considerations
- instead of the cache library, simply redis get and set are used.
- for limiting the vote of users, we have UserVoteLimiter which keeps track of user daily votes.
- there is a chance that user may vote for more than 100 polls, if use the script that can request 
so many votes at the same time, if we want to avoid that, we should use locking mechanism (not implemented).
- there is stress_test script using wrk,
to run it you should have wrk and jq installed on your machine. then run ```stresstest.sh```
- no cache is used for the stats api, because the rate of voting causes to invalidate cache almost instantly,
unless we don't care about the actual count of selected options. index on db would be enough for now.


### future growth
- I believe we cannot use cache for the route of getting polls list, as it depends on each user and page, the volume of
data is huge, and cache invalidation is so hard, when there is a new poll.
maybe we can have api to return id of last 1000 polls that user skips or voted, in this case,
we can send these ids to application. then we can cache the list of polls for all users,
and in the application, the app could check if a poll already voted or skipped by using the former 
api and remove them from the list.
