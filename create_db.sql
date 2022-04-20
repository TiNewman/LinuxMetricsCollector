/*
Author: Titan Newman
Date: 4/19/2022

Creation Script for the MetricsCollectorDB.
*/

DROP DATABASE IF EXISTS MetricsCollectorDB;

GO

CREATE DATABASE MetricsCollectorDB;

GO

USE MetricsCollectorDB;

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
	-- For now this arent used as we are only working with Process and CPU.
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

	CONSTRAINT pk_process_processID_collectorID_PID PRIMARY KEY (processID, collectorID, PID),
	CONSTRAINT fk_process_collector_collectorID FOREIGN KEY (collectorID) REFERENCES COLLECTOR(collectorID)
);

-- HISTORY TABLES

CREATE TABLE CPU_AVERAGE (
	cpuAverageID BIGINT NOT NULL IDENTITY(0,1),
	averageUsage FLOAT NOT NULL,
	
	CONSTRAINT pk_cpuaverage_cpuaverageID PRIMARY KEY (cpuAverageID)
);

CREATE TABLE MEMORY_AVERAGE (
	memoryAverageID BIGINT NOT NULL IDENTITY(0,1),
	averageUsage FLOAT NOT NULL,
	averageAvailability FLOAT NOT NULL,

	CONSTRAINT pk_memoryaverage_memoryaverageID PRIMARY KEY (memoryAverageID)
);

CREATE TABLE DISK_AVERAGE (
	diskAverageID BIGINT NOT NULL IDENTITY(0,1),
  averageUsage FLOAT NOT NULL,
	averageAvailability FLOAT NOT NULL,

	CONSTRAINT pk_diskaverage_diskaverageID PRIMARY KEY (diskAverageID)
);

CREATE TABLE COLLECTOR_HISTORY (
	collectorHistoryID BIGINT NOT NULL IDENTITY(0,1),
	timeCollectedStart DATE NOT NULL,
	timeCollectedEnd DATE NOT NULL,
	-- For now this arent used as we are only working with Process and CPU.
	cpuAverageID BIGINT, -- NOT NULL
	memoryAverageID BIGINT, -- NOT NULL
	diskAverageID BIGINT, -- NOT NULL

	CONSTRAINT pk_collectorhistory_collectorHistoryID PRIMARY KEY (collectorHistoryID),
	CONSTRAINT fk_collectorhistory_cpuaverage_averagecpuID FOREIGN KEY (cpuAverageID) REFERENCES CPU_AVERAGE(cpuAverageID),
	CONSTRAINT fk_collectorhistory_memoryaverage_averagememoryID FOREIGN KEY (memoryAverageID) REFERENCES MEMORY_AVERAGE(memoryAverageID),
	CONSTRAINT fk_collectorhistory_diskaverage_averagediskID FOREIGN KEY (diskAverageID) REFERENCES DISK_AVERAGE(diskAverageID)
);

CREATE TABLE PROCESS_HISTORY (
	processHistoryID BIGINT NOT NULL IDENTITY(0,1),
	collectorHistoryID BIGINT NOT NULL,
	PID BIGINT NOT NULL,
	name VARCHAR(100),
	status VARCHAR(20) NOT NULL,
	cpuUsage FLOAT,
	memoryUsage FLOAT,
	diskUsage FLOAT,
	executionTime FLOAT,

	CONSTRAINT pk_processhistory_processhistoryID_collectorhistoryID PRIMARY KEY (processHistoryID, collectorHistoryID),
	CONSTRAINT fk_processhistory_collectorhistory_collectorhistoryID FOREIGN KEY (collectorHistoryID) REFERENCES COLLECTOR_HISTORY(collectorHistoryID)
);

GO


-- Stored Procedure for Purging data
/*

CREATE PROCEDURE PurgeData
AS

DECLARE @startDate DATE = CONVERT(DATE, (SELECT TOP 1 timeCollected FROM COLLECTOR ORDER BY timeCollected ASC));
DECLARE @endDate DATE = CONVERT(DATE, (SELECT DATEADD(DD, 3, @startDate)));

--SELECT @startDate AS "Start time", @endDate AS "NEWTIME";
--SELECT timeCollected FROM COLLECTOR ORDER BY timeCollected ASC;
DECLARE @endCollectorID BIGINT = (SELECT TOP 1 collectorID FROM COLLECTOR WHERE CONVERT(DATE, timeCollected) = @endDate ORDER BY timeCollected DESC);

DECLARE @cpuUsage FLOAT = 0.0;
DECLARE @diskUsage FLOAT = 0.0;
DECLARE @memoryUsage FLOAT = 0.0;
DECLARE @diskAvailability FLOAT = 0.0;
DECLARE @memoryAvailability FLOAT = 0.0;

DECLARE @startingID BIGINT = (SELECT TOP 1 cpuID FROM CPU ORDER BY cpuID ASC);
DECLARE @endID BIGINT = (SELECT TOP 1 cpuID FROM CPU WHERE cpuID IN (SELECT cpuID FROM COLLECTOR WHERE collectorID = @endCollectorID) ORDER BY cpuID DESC);

DECLARE @count FLOAT = (@endID - @startingID);

IF @count >= 1
BEGIN

	IF @startingID = 0 
	BEGIN
		SET @count = @count + 1;
	END


	IF @startingID = (SELECT TOP 1 diskID FROM DISK ORDER BY diskID ASC) AND (@startingID = (SELECT TOP 1 memoryID FROM MEMORY ORDER BY memoryID ASC)) 
	BEGIN

		WHILE @startingID <= @endID
		BEGIN
			-- Usages
			SET @cpuUsage = @cpuUsage + (SELECT TOP 1 usage FROM CPU WHERE cpuID = @startingID);
			SET @diskUsage = @diskUsage + (SELECT TOP 1 usage FROM DISK WHERE diskID = @startingID);
			--SELECT TOP 1 usage FROM MEMORY WHERE memoryID = @startingID;
			SET @memoryUsage = @memoryUsage + (SELECT TOP 1 usage FROM MEMORY WHERE memoryID = @startingID);
			-- Availabilities
			SET @diskAvailability = @diskAvailability + (SELECT TOP 1 availability FROM DISK WHERE diskID = @startingID);
			--SELECT TOP 1 availability FROM MEMORY WHERE memoryID = @startingID;
			SET @memoryAvailability = @memoryAvailability + (SELECT TOP 1 availability FROM MEMORY WHERE memoryID = @startingID);

			-- Increase ID pointer.
			SET @startingID = @startingID + 1;
		END
	
		-- Getting the average, but dividing by 3 as per day.
		SET @cpuUsage = CAST(ROUND((@cpuUsage / CAST(@count AS FLOAT)) / CAST(3 AS FLOAT), 2) AS NUMERIC(36,2));
		SET @diskUsage = CAST(ROUND((@diskUsage / CAST(@count AS FLOAT)) / CAST(3 AS FLOAT), 2) AS NUMERIC(36,2));
		SET @memoryUsage = CAST(ROUND((@memoryUsage / CAST(@count AS FLOAT)) / CAST(3 AS FLOAT), 2) AS NUMERIC(36,2));
		SET @diskAvailability = CAST(ROUND((@diskAvailability / CAST(@count AS FLOAT)) / CAST(3 AS FLOAT), 2) AS NUMERIC(36,2));
		SET @memoryAvailability = CAST(ROUND((@memoryAvailability / CAST(@count AS FLOAT)) / CAST(3 AS FLOAT), 2) AS NUMERIC(36,2));

		SELECT @cpuUsage AS "CPU-USAGE", @diskUsage AS "DISK-USAGE", @memoryUsage AS "MEMORY-USAGE", @diskAvailability AS "DISK-AVA", @memoryAvailability AS "MEMORY-AVA" ;

	END

	ELSE
	BEGIN

		PRINT N'Error: Purge Stored_Procedure_CPU_DISK_MEMORY --> Inconsistent data.';
	END
END

ELSE
BEGIN

	PRINT N'Error: Purge Stored_Procedure_CPU_DISK_MEMORY --> Count for cycling through CPU/DISK/MEMORY was negative.';
END

-- INSERT into CPU/DISK/MEMORY _AVERAGE tables
INSERT INTO CPU_AVERAGE VALUES (@cpuUsage);
DECLARE @insertedCpuID BIGINT = (SELECT cpuAverageID FROM CPU_AVERAGE WHERE cpuAverageID = SCOPE_IDENTITY() AND averageUsage = @cpuUsage);
INSERT INTO DISK_AVERAGE VALUES (@diskUsage, @diskAvailability);
DECLARE @insertedDiskID BIGINT = (SELECT diskAverageID FROM DISK_AVERAGE WHERE diskAverageID = SCOPE_IDENTITY() AND averageUsage = @diskUsage AND averageAvailability = @diskAvailability);
INSERT INTO MEMORY_AVERAGE VALUES (@memoryUsage, @memoryAvailability);
DECLARE @insertedMemoryID BIGINT = (SELECT memoryAverageID FROM MEMORY_AVERAGE WHERE memoryAverageID = SCOPE_IDENTITY() AND averageUsage = @memoryUsage AND averageAvailability = @memoryAvailability);


-- INSERT into COLLECTOR_HISTORY table
INSERT INTO COLLECTOR_HISTORY VALUES (@startDate, @endDate, @insertedCpuID, @insertedMemoryID, @insertedDiskID);




--SET @startingID = (SELECT TOP 1 processID FROM PROCESS ORDER BY processID ASC);
--SET @endID = (SELECT TOP 1 cpuID FROM CPU ORDER BY cpuID DESC);
	
GO;


*/