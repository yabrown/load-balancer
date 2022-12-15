Things to test:
distribution of load among servers-- server fairness
average response time-- client optimization
response time distribution-- client fairness
processing time average-- might be willing to wait for data, as long as done quickly (ex. video loading, stock trading stuff)


Test Cases:
equal_cores_best_dist-- 5 servers of 1 core each, with a purposefully even load distribution 
equal_cores_worst_dist-- 5 servers of 1 core each, with a purposefully bad load distribution 
equal_cores_random_dist-- 5 servers of 1 core each, with random distribution 