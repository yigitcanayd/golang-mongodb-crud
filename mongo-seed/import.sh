#! /bin/bash

mongoimport --host mongodb --db golang --collection people --type json --file /mongo-seed/people.json --jsonArray