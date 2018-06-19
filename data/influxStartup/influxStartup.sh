#!/bin/bash
influx -execute 'CREATE RETENTION POLICY "sixMonths" ON "doppler" DURATION 26w REPLICATION 1 DEFAULT'