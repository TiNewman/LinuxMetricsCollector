/*
Author: Titan Newman
Date: 2/24/2022

Creation Script for the MetricsCollectorDB.
*/

CREATE DATABASE MetricsCollectorDB;

GO

USE MetricsCollectorDB;

GO

CREATE TABLE CPU (
	cpuID BIGINT NOT NULL IDENTITY(0,1),
    usage FLOAT NOT NULL,
	availability FLOAT NOT NULL,

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
    cpuID BIGINT NOT NULL,
	memoryID BIGINT NOT NULL,
	diskID BIGINT NOT NULL,

	CONSTRAINT pk_collector_collectorID PRIMARY KEY (collectorID),
	CONSTRAINT fk_collector_cpu_cpuID FOREIGN KEY (cpuID) REFERENCES CPU(cpuID),
	CONSTRAINT fk_collector_memory_memoryID FOREIGN KEY (memoryID) REFERENCES MEMORY(memoryID),
	CONSTRAINT fk_collector_disk_diskID FOREIGN KEY (diskID) REFERENCES DISK(diskID)
);

CREATE TABLE PROCESS (
	processID BIGINT NOT NULL IDENTITY(0,1),
	collectorID BIGINT NOT NULL,
	name VARCHAR(100) NOT NULL,
	status VARCHAR(20) NOT NULL,
	cpuUsage FLOAT,
	memoryUsage FLOAT,
	diskUsage FLOAT,
	executionTime FLOAT,

	CONSTRAINT pk_process_processID PRIMARY KEY (processID),
	CONSTRAINT fk_process_collector_collectorID FOREIGN KEY (collectorID) REFERENCES COLLECTOR(collectorID)
);

GO