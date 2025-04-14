SELECT
    Name,Dept
FROM
    employee
    INNER JOIN department ON employee.DEPT_ID = department.ID
