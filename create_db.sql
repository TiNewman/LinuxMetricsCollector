/*
Author: Titan Newman
Date: 2/28/2022

Creation Script for the MetricsCollectorDB.
*/

/*
If partitioning doesnt work:
Truncate Table: https://docs.microsoft.com/en-us/sql/t-sql/statements/truncate-table-transact-sql?view=sql-server-ver15




*/

DROP DATABASE IF EXISTS MetricsCollectorDB;

GO

CREATE DATABASE MetricsCollectorDB;

GO

USE MetricsCollectorDB;

GO

CREATE TABLE CPU (
	cpuID BIGINT NOT NULL IDENTITY(0,1),
  usage FLOAT NOT NULL,

	CONSTRAINT pk_cpu_cpuID PRIMARY KEY (cpuID)
);

CREATE TABLE MEMORY (
	memoryID BIGINT NOT NULL IDENTITY(0,1),
  usage FLOAT NOT NULL,
	availability FLOAT NOT NULL,

	CONSTRAINT pk_memory_memoryID PRIMARY KEY (memoryID)
);

CREATE TABLE DISK (
	diskID BIGINT NOT NULL IDENTITY(0,1),
  usage FLOAT NOT NULL,
	availability FLOAT NOT NULL,

	CONSTRAINT pk_disk_diskID PRIMARY KEY (diskID)
);

CREATE TABLE COLLECTOR (
	collectorID BIGINT NOT NULL IDENTITY(0,1),
	timeCollected DATETIME2 NOT NULL, -- Pay attention to how the data needs to be formatted here!
	-- For now this arent used as we are only working with Process.
  cpuID BIGINT, -- NOT NULL
	memoryID BIGINT, -- NOT NULL
	diskID BIGINT, -- NOT NULL

	CONSTRAINT pk_collector_collectorID PRIMARY KEY (collectorID),
	CONSTRAINT fk_collector_cpu_cpuID FOREIGN KEY (cpuID) REFERENCES CPU(cpuID),
	CONSTRAINT fk_collector_memory_memoryID FOREIGN KEY (memoryID) REFERENCES MEMORY(memoryID),
	CONSTRAINT fk_collector_disk_diskID FOREIGN KEY (diskID) REFERENCES DISK(diskID)
);

CREATE TABLE PROCESS (
	processID BIGINT NOT NULL IDENTITY(0,1),
	collectorID BIGINT NOT NULL,
	PID BIGINT NOT NULL, -- This is the actual PID of the process from PROCFS.
	name VARCHAR(100),
	status VARCHAR(20) NOT NULL,
	cpuUsage FLOAT,
	memoryUsage FLOAT,
	diskUsage FLOAT,
	executionTime FLOAT,

	CONSTRAINT pk_process_processID PRIMARY KEY (processID),
	CONSTRAINT fk_process_collector_collectorID FOREIGN KEY (collectorID) REFERENCES COLLECTOR(collectorID)
);

GO

/*
DECLARE @dateHolder date = '12-31-9999';
DECLARE @datetime2Holder datetime2 = @dateHolder;

CREATE PARTITION FUNCTION dateRangePartitionFunc (datetime)
    AS RANGE FOR VALUES (DATEDIFF(day , SYSDATETIME() , @datetime2Holder)) ;  
GO

CREATE PARTITION SCHEME dateRangePartitionSchme 
    AS PARTITION dateRangePartitionFunc  
    ALL TO ('PRIMARY');
GO  

CREATE TABLE COLLECTORXX (
	collectorID BIGINT NOT NULL IDENTITY(0,1),
	timeCollected DATETIME NOT NULL,
	cpuID BIGINT, -- NOT NULL
	memoryID BIGINT, -- NOT NULL
	diskID BIGINT, -- NOT NULL

	--CONSTRAINT pk_collector_collectorID2 PRIMARY KEY (collectorID), primary key has to be timedate or move this back to a sub-tree of datetime 
) 
ON dateRangePartitionSchme (timeCollected);  
GO
*/