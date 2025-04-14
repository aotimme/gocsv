# [get sql query string from file](https://github.com/aotimme/gocsv/issues/79)

On Windows when using cmd it will not allow one to run quoted strings over multiple lines. The sql query strings can be very long making it hard to read. Would be nice if the sql subcommand could get the sql query string from a file like this, say:

```
gocsv sql -f mycalc.sql input.csv
```

## My analysis

OP wants to do something like:

```none
-- employee.csv --
ID, Name,    Age, Salary,  Dept_ID
 1, Abiola,   29, 70000.0,       2
 2, Bolade,   25, 68000.0,       3
 3, Chima,    23, 66000.0,       2
 4, Dinya,    25, 80000.0,       1
 5, Ekon,     27, 85000.0,       1
 6, Feranmi,  22, 65000.0,       2
 7, Ifama,    24, 65000.0,       2

-- department.csv --
ID, Dept
 1, DevOps
 2, Engineering
 3, Finance

-- join.sql --
SELECT
    Name,Dept
FROM
    employee
    INNER JOIN department ON employee.DEPT_ID = department.ID

```

```none
gocsv sql -f=join.sql department.csv employee.csv
```

and get:

```none
Name,    Dept
Abiola,  Engineering
Bolade,  Finance
Chima,   Engineering
Dinya,   DevOps
Ekon,    DevOps
Feranmi, Engineering
Ifama,   Engineering
```

## My response
