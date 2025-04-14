#!/bin/sh

cd Issue-79

go run .. sql -f=join.sql department.csv employee.csv >/dev/null || exit 1

go run .. sql -q='SELECT SUM(Salary) AS Total FROM employee' -f=FooBarBaz employee.csv >/dev/null 2>&1
if [ $? -ne 1 ]; then
    echo 'expected error about using both -q and -f flags together'
    exit 1
fi
