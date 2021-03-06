/*
Team: Titan, Jordan, and James 
Author: Titan Newman
Date: 4/29/2022

Creation Script for the MetricsCollectorDB.
This includes tables and a single stored procedure.
	Stored procedure averages data and moves it to history tables.
	Runs to insure Atomicity and Consistency for data being moved
	(uses transactions).
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
	size FLOAT NOT NULL,

	CONSTRAINT pk_memory_memoryID PRIMARY KEY (memoryID)
);

CREATE TABLE DISK (
	diskID BIGINT NOT NULL IDENTITY(0,1),
  usage FLOAT NOT NULL,
	size FLOAT NOT NULL,

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
	averageSize FLOAT NOT NULL,

	CONSTRAINT pk_memoryaverage_memoryaverageID PRIMARY KEY (memoryAverageID)
);

CREATE TABLE DISK_AVERAGE (
	diskAverageID BIGINT NOT NULL IDENTITY(0,1),
  averageUsage FLOAT NOT NULL,
	averageSize FLOAT NOT NULL,

	CONSTRAINT pk_diskaverage_diskaverageID PRIMARY KEY (diskAverageID)
);

CREATE TABLE COLLECTOR_HISTORY (
	collectorHistoryID BIGINT NOT NULL IDENTITY(0,1),
	timeCollectedStart DATE NOT NULL,
	timeCollectedEnd DATE NOT NULL,
	cpuAverageID BIGINT NOT NULL,
	memoryAverageID BIGINT NOT NULL,
	diskAverageID BIGINT NOT NULL,

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


-- Stored Procedure for Purging data + averaging.
IF (OBJECT_ID('PRC_PurgeData') IS NOT NULL)
BEGIN
	DROP PROCEDURE PRC_PurgeData;
END

GO

CREATE PROCEDURE PRC_PurgeData 
AS

DECLARE @startDate DATE = CONVERT(DATE, (SELECT TOP 1 timeCollected FROM COLLECTOR ORDER BY timeCollected ASC));
DECLARE @endDate DATE = CONVERT(DATE, (SELECT DATEADD(DD, 3, @startDate)));
DECLARE @currentDate DATE = CONVERT(DATE, (SELECT TOP 1 timeCollected FROM COLLECTOR ORDER BY timeCollected DESC));

-- If we have gone past 3 days, purge the database.
IF (SELECT DATEDIFF(DAY, @startDate, @currentDate) AS DateDiff) >= 3 
BEGIN
	
	DECLARE @endCollectorID BIGINT = (SELECT TOP 1 collectorID FROM COLLECTOR WHERE CONVERT(DATE, timeCollected) = @endDate ORDER BY timeCollected DESC);

	DECLARE @cpuUsage FLOAT = 0.0;
	DECLARE @diskUsage FLOAT = 0.0;
	DECLARE @memoryUsage FLOAT = 0.0;
	DECLARE @diskSize FLOAT = 0.0;
	DECLARE @memorySize FLOAT = 0.0;

	DECLARE @startingID BIGINT = (SELECT TOP 1 cpuID FROM CPU ORDER BY cpuID ASC);
	DECLARE @endID BIGINT = (SELECT TOP 1 cpuID FROM CPU WHERE cpuID IN (SELECT cpuID FROM COLLECTOR WHERE collectorID = @endCollectorID) ORDER BY cpuID DESC);

	DECLARE @count FLOAT = (@endID - @startingID);

	IF (@count >= 1)
	BEGIN
		
		IF (@startingID = 0) 
		BEGIN
			SET @count = @count + 1;
		END


		IF (@startingID = (SELECT TOP 1 diskID FROM DISK ORDER BY diskID ASC) AND (@startingID = (SELECT TOP 1 memoryID FROM MEMORY ORDER BY memoryID ASC))) 
		BEGIN

			WHILE (@startingID <= @endID)
			BEGIN
				-- Usages
				SET @cpuUsage = @cpuUsage + (SELECT TOP 1 usage FROM CPU WHERE cpuID = @startingID);
				SET @diskUsage = @diskUsage + (SELECT TOP 1 usage FROM DISK WHERE diskID = @startingID);
				SET @memoryUsage = @memoryUsage + (SELECT TOP 1 usage FROM MEMORY WHERE memoryID = @startingID);
				-- Availabilities
				SET @diskSize = @diskSize + (SELECT TOP 1 size FROM DISK WHERE diskID = @startingID);
				SET @memorySize = @memorySize + (SELECT TOP 1 size FROM MEMORY WHERE memoryID = @startingID);

				-- Increase ID pointer.
				SET @startingID = @startingID + 1;
			END
	
			-- Getting the average, but dividing by 3 as per day.
			SET @cpuUsage = CAST(ROUND((@cpuUsage / CAST(@count AS FLOAT)) / CAST(3 AS FLOAT), 2) AS NUMERIC(36,2));
			SET @diskUsage = CAST(ROUND((@diskUsage / CAST(@count AS FLOAT)) / CAST(3 AS FLOAT), 2) AS NUMERIC(36,2));
			SET @memoryUsage = CAST(ROUND((@memoryUsage / CAST(@count AS FLOAT)) / CAST(3 AS FLOAT), 2) AS NUMERIC(36,2));
			SET @diskSize = CAST(ROUND((@diskSize / CAST(@count AS FLOAT)) / CAST(3 AS FLOAT), 2) AS NUMERIC(36,2));
			SET @memorySize = CAST(ROUND((@memorySize / CAST(@count AS FLOAT)) / CAST(3 AS FLOAT), 2) AS NUMERIC(36,2));

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

	BEGIN TRANSACTION

	-- INSERT into CPU/DISK/MEMORY _AVERAGE tables
	INSERT INTO CPU_AVERAGE VALUES (@cpuUsage);
	DECLARE @insertedCpuID BIGINT = (SELECT cpuAverageID FROM CPU_AVERAGE WHERE cpuAverageID = SCOPE_IDENTITY());
	
	IF (@@ERROR <> 0)
	BEGIN
		RAISERROR ('Error inserting into CPU_AVERAGE',-1,-1);
		ROLLBACK TRANSACTION;
		RETURN;
	END

	INSERT INTO DISK_AVERAGE VALUES (@diskUsage, @diskUsage);
	DECLARE @insertedDiskID BIGINT = (SELECT diskAverageID FROM DISK_AVERAGE WHERE diskAverageID = SCOPE_IDENTITY());

	IF (@@ERROR <> 0)
	BEGIN
		RAISERROR ('Error inserting into DISK_AVERAGE',-1,-1);
		ROLLBACK TRANSACTION;
		RETURN;
	END
	
	INSERT INTO MEMORY_AVERAGE VALUES (@memoryUsage, @memoryUsage);
	DECLARE @insertedMemoryID BIGINT = (SELECT memoryAverageID FROM MEMORY_AVERAGE WHERE memoryAverageID = SCOPE_IDENTITY());

	IF (@@ERROR <> 0)
	BEGIN
		RAISERROR ('Error inserting into MEMORY_AVERAGE',-1,-1);
		ROLLBACK TRANSACTION;
		RETURN;
	END


	-- INSERT into COLLECTOR_HISTORY table
	INSERT INTO COLLECTOR_HISTORY VALUES (@startDate, @endDate, @insertedCpuID, @insertedMemoryID, @insertedDiskID);
	DECLARE @collectorHistoryID BIGINT = (SELECT TOP 1 collectorHistoryID FROM COLLECTOR_HISTORY ORDER BY collectorHistoryID DESC);

	IF (@@ERROR <> 0)
	BEGIN
		RAISERROR ('Error inserting into COLLECTION_HISTORY',-1,-1);
		ROLLBACK TRANSACTION;
		RETURN;
	END

	-- PROCCESS -> PROCESS_HISTORY
	SET @startingID = (SELECT TOP 1 processID FROM PROCESS ORDER BY processID ASC);
	SET @endID = (SELECT TOP 1 processID FROM PROCESS WHERE collectorID = @endCollectorID ORDER BY processID DESC);
	
	IF @endID IS NULL 
	BEGIN
		SET @endID = (SELECT TOP 1 processID FROM PROCESS WHERE collectorID < @endCollectorID ORDER BY processID DESC);
	END 

	SET @count = (@endID - @startingID);

	IF (@count >= 1)
	BEGIN

		IF (@startingID = 0) 
		BEGIN
			SET @count = @count + 1;
		END

		WHILE (@startingID <= @endID)
		BEGIN

			INSERT INTO PROCESS_HISTORY (collectorHistoryID, PID, name, status, cpuUsage, memoryUsage, diskUsage, executionTime)
			SELECT @collectorHistoryID, PID, name, status, cpuUsage, memoryUsage, diskUsage, executionTime 
			FROM PROCESS WHERE processID = @startingID;

			IF (@@ERROR <> 0)
			BEGIN
				RAISERROR ('Error inserting into PROCESS_HISTORY',-1,-1);
				ROLLBACK TRANSACTION;
				RETURN;
			END

			-- DELETE from PROCESS
			DELETE FROM PROCESS WHERE processID = @startingID;

			IF (@@ERROR <> 0)
			BEGIN
				RAISERROR ('Error deleting from PROCESS',-1,-1);
				ROLLBACK TRANSACTION;
				RETURN;
			END

			-- Increase ID
			SET @startingID = @startingID + 1;

		END
	END

	ELSE
	BEGIN

		PRINT N'Purge Stored_Procedure_PROCESS --> Count for cycling through PROCESS was negative.';
	END


	-- DELETE from CPU/DISK/MEMORY and COLLECTOR

	DECLARE @startCollectorID BIGINT = (SELECT TOP 1 collectorID FROM COLLECTOR ORDER BY timeCollected ASC);

	DECLARE @cpuIdDelete BIGINT;
	DECLARE @diskIdDelete BIGINT;
	DECLARE @memoryIdDelete BIGINT;

	WHILE (@startCollectorID <= @endCollectorID) 
	BEGIN

		SET @cpuIdDelete = (SELECT cpuID FROM COLLECTOR WHERE collectorID = @startCollectorID);
		SET @diskIdDelete = (SELECT diskID FROM COLLECTOR WHERE collectorID = @startCollectorID);
		SET @memoryIdDelete = (SELECT memoryID FROM COLLECTOR WHERE collectorID = @startCollectorID);
		
		--DELETE from tables.
		DELETE FROM COLLECTOR WHERE collectorID = @startCollectorID;
		IF (@@ERROR <> 0)
		BEGIN
			RAISERROR ('Error deleting from COLLECTOR',-1,-1);
			ROLLBACK TRANSACTION;
			RETURN;
		END

		DELETE FROM CPU WHERE cpuID = @cpuIdDelete;

		IF (@@ERROR <> 0)
		BEGIN
			RAISERROR ('Error deleting from CPU',-1,-1);
			ROLLBACK TRANSACTION;
			RETURN;
		END

		DELETE FROM DISK WHERE diskID = @diskIdDelete;

		IF (@@ERROR <> 0)
		BEGIN
			RAISERROR ('Error deleting from DISK',-1,-1);
			ROLLBACK TRANSACTION;
			RETURN;
		END

		DELETE FROM MEMORY WHERE memoryID = @memoryIdDelete;

		IF (@@ERROR <> 0)
		BEGIN
			RAISERROR ('Error deleting from MEMORY',-1,-1);
			ROLLBACK TRANSACTION;
			RETURN;
		END

		-- Increase ID.
		SET @startCollectorID = @startCollectorID + 1;

	END

	COMMIT TRANSACTION;

END
	
GO
