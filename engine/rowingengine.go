package engine

import (
	"math/rand"
	"time"
)

const (
	CENTISECONDS_PER_SECOND = 100
	MILLISECONDS_PER_CENTISECOND = 10
	DECIMETERS_PER_METER = 10
	TWO_MINUTE_SPLIT_SPEED = 4.16  // as meters per second
)

const (
	second		Size = 0
	half		Size = 1
	quarter		Size = 2
	tenth		Size = 3
)

type rowingEngine struct {
	startTime time.Time
	statusRate uint8
}

func NewRowingEngine() rowingEngine {
	return rowingEngine{time.Now(), half}
}

/*-------- Private Methods ---------*/

/*
	Converts a decimal value to a little-endian byte array. This is 
	useful for converting values to the [lo, hi] or [lo, mid, hi] 
	format used in PM5 characteristics.
	E.g. decimalToHexBytes(12000, 2) returns [224, 46]
	because 224 + (46 * 256) = 12000
*/
func decimalToBytes(value int64, len int) []byte {	
	bytes := make([]byte, len, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(value % 256)
		value = value / int64(256)	
	}
	return bytes
}

/*
	Generates the lo, mid, and high bytes for the time elapsed since the
	rowingEngine struct was initialized. LSB refers to 0.01S. 
*/
func (eng rowingEngine) generateElapsedTimeBytes() []byte {
	elapsedDuration := time.Since(eng.startTime).Milliseconds() / int64(MILLISECONDS_PER_CENTISECOND)
	return decimalToBytes(elapsedDuration, 3)
}

/*
	Generates the lo, mid, hi bytes for distance rowed. LSB refers to 0.1m.
	This is simplified to distance that would have been rowed since the rowingEngine
	was initialized, assuming a 2:00 split. 
*/
func (eng rowingEngine) generateDistanceBytes() []byte {
	elapsedSeconds := time.Since(eng.startTime).Seconds() 
	distance := elapsedSeconds * TWO_MINUTE_SPLIT_SPEED * DECIMETERS_PER_METER
	return decimalToBytes(int64(distance), 3)
}

/*
	Generates the lo and hi bytes for a realistic "Split/Int Avg Pace" value
	For now "realistic" is defined as a random value between 120s and 130s,
	corresponding to a 2:00-2:10 split. 
*/
func (eng rowingEngine) generateSplitBytes() []byte {
	splitCentiseconds := (120 + rand.Intn(10)) * CENTISECONDS_PER_SECOND
	return decimalToBytes(int64(splitCentiseconds), 2)
}

/*
	Generates a realistic stroke rate byte value.
	For now "realistic" is defined as a random value between 25 and 30. 
*/
func (eng rowingEngine) generateStrokeRateBytes() []byte {
	strokeRate := 25 + rand.Intn(5)
	return decimalToBytes(int64(strokeRate), 1)
}

/*-------- Public Methods ---------*/

/*
	returns a number of milliseconds to wait between each status packet
	in order to abide by the configured statusRate value. 
	This corresponds to the possible values of characteristic 0x0034 from the
	concept2 bluetooth specification
	0 – 1 sec
	1 - 500ms (default if characteristic is not explicitly set by the
	app)
	2 – 250ms
	3 – 100ms
*/
func (eng rowingEngine) GetStatusDelay() []byte {
	switch eng.statusRate {
		case 0: return 1000
		case 1: return 500
		case 2: return 250
		case 3: return 100
		default: 500
	}
}

//TODO: add a method to allow the statusDelay to be set from a bluetooth characteristic

/*
	Generates a 'General Status' characteristic of 19 bytes.
	Only the bytes labeled with a * are populated with non-zero values:
	0 - Elapsed Time Lo (0.01 sec lsb), 		*
	1 - Elapsed Time Mid,										*
	2 - Elapsed Time High,									*
	3 - Distance Lo (0.1 m lsb), 						*
	4 - Distance Mid,												*
	5 - Distance High,											*
	6 - Workout Type (enum), 
	7 - Interval Type
	8 - Workout State (enum), 
	9 - Rowing State (enum), 
	10 - Stroke State (enum), 
	11 - Total Work Distance Lo, 
	12 - Total Work Distance Mid,
	13 - Total Work Distance Hi,
	14 - Workout Duration Lo (if time, 0.01 sec lsb), 
	15 - Workout Duration Mid,
	16 - Workout Duration Hi,
	17 - Workout Duration Type (enum),
	18 - Drag Factor 
*/
func (eng rowingEngine) GenerateGeneralStatusChar() []byte {
	bytes := make([]byte, 19, 19)

	copy(bytes[0:], eng.generateElapsedTimeBytes()) 
	copy(bytes[3:], eng.generateDistanceBytes())
	
	return bytes
}

/*
	Generates an 'Additional Status 1' characteristic of 17 bytes.
	Only the bytes labeled with a * are populated with non-zero values:
		0 - Elapsed Time Lo (0.01 sec lsb), 		*
		1 - Elapsed Time Mid,										*
		2 - Elapsed Time High,									*
		3 - Speed Lo (0.001m/s lsb), 
		4 - Speed Hi,
		5 - Stroke Rate (strokes/min), 					*
		6 - Heartrate (bpm, 255=invalid), 
		7 - Current Pace Lo (0.01 sec lsb),
		8 - Current Pace Hi,
		9 - Average Pace Lo (0.01 sec lsb),
		10 - Average Pace Hi,
		11 - Rest Distance Lo, 
		12 - Rest Distance Hi,
		13 - Rest Time Lo, (0.01 sec lsb) 
		14 - Rest Time Mid,
		15 - Rest Time Hi
		16 - Erg Machine Type
*/
func (eng rowingEngine) GenerateAdditionalStatus1Char() []byte {
	bytes := make([]byte, 17, 17)

	copy(bytes[0:], eng.generateElapsedTimeBytes())
	copy(bytes[5:], eng.generateStrokeRateBytes())
	
	return bytes
}

/*
	Generates an 'Additional Status 2' characteristic of 20 bytes.
	Only the bytes labeled with a * are populated with non-zero values:
		0 - Elapsed Time Lo (0.01 sec lsb), 		*
		1 - Elapsed Time Mid,										*
		2 - Elapsed Time High,									*
		3 - Interval Count, 
		4 - Average Power Lo, 
		5 - Average Power Hi,
		6 - Total Calories Lo (cals), 
		7 - Total Calories Hi,
		8 - Split/Int Avg Pace Lo (0.01 sec lsb), 	*
		9 - Split/Int Avg Pace Hi,									*
		10 - Split/Int Avg Power Lo (watts), 
		11 - Split/Int Avg Power Hi,
		12 - Split/Int Avg Calories Lo (cals/hr), 
		13 - Split/Interval Avg Calories Hi,
		14 - Last Split Time Lo (0.1 sec lsb), 
		15 - Last Split Time Mid,
		16 - Last Split Time High,
		17 - Last Split Distance Lo, 
		18 - Last Split Distance Mid,
		19 - Last Split Distance Hi
*/
func (eng rowingEngine) GenerateAdditionalStatus2Char() []byte {
	bytes := make([]byte, 20, 20)

	copy(bytes[0:], eng.generateElapsedTimeBytes())
	copy(bytes[8:], eng.generateSplitBytes())
	
	return bytes
}

/*
	Generates a 'Stroke Data' characteristic of 20 bytes.
	Only the bytes labeled with a * are populated with non-zero values:
		0 - Elapsed Time Lo (0.01 sec lsb), *
		1 - Elapsed Time Mid, *
		2 - Elapsed Time High, *
		3 - Distance Lo (0.1 m lsb),
		4 - Distance Mid,
		5 - Distance High,
		6 - Drive Length (0.01 meters, max = 2.55m),
		7 - Drive Time (0.01 sec, max = 2.55 sec),
		8 - Stroke Recovery Time Lo (0.01 sec, max = 655.35 sec),
		9 - Stroke Recovery Time Hi
		10 - Stroke Distance Lo (0.01 m, max=655.35m)
		11 - Stroke Distance Hi,
		12 - Peak Drive Force Lo (0.1 lbs of force, max=6553.5m)
		13 - Peak Drive Force Hi,
		14 - Average Drive Force Lo (0.1 lbs of force, max=6553.5m),
		15 - Average Drive Force Hi,
		16 - Work Per Stroke Lo (0.1 Joules, max=6553.5 Joules),
		17 - Work Per Stroke Hi
		18 - Stroke Count Lo,
		19 - Stroke Count Hi,
*/
func (eng rowingEngine) GenerateStrokeDataChar() []byte {
	bytes := make([]byte, 20, 20)

	copy(bytes[0:], eng.generateElapsedTimeBytes())
	// copy(bytes[8:], eng.generateSplitBytes())
	
	return bytes
}

/*
	Generates a 'Stroke Data 2' characteristic of 15 bytes.
	Only the bytes labeled with a * are populated with non-zero values:
		0 - Elapsed Time Lo (0.01 sec lsb), *
		1 - Elapsed Time Mid, *
		2 - Elapsed Time High, *
		3 - Stroke Power Lo (watts),
		4 - Stroke Power Hi,
		5 - Stroke Calories Lo (cal/hr),
		6 - Stroke Calories Hi,
		7 - Stroke Count Lo, 
		8 - Stroke Count Hi,
		9 - Projected Work Time Lo (secs),
		10 - Projected Work Time Mid,
		11 - Projected Work Time Hi,
		12 - Projected Work Distance Lo (meters),
		13 - Projected Work Distance Mid,
		14 - Projected Work Distance Hi		
*/
func (eng rowingEngine) GenerateStrokeData2Char() []byte {
	bytes := make([]byte, 15, 15)

	copy(bytes[0:], eng.generateElapsedTimeBytes())
	return bytes
}

/*
	Generates a 'Split/Interval' characteristic of 18 bytes.
	Only the bytes labeled with a * are populated with non-zero values:
		0 - Elapsed Time Lo (0.01 sec lsb), *
		1 - Elapsed Time Mid, *
		2 - Elapsed Time High, *
		3 - Distance Lo (0.1 m lsb),
		4 - Distance Mid,
		5 - Distance High,
		6 - Split/Interval Time Lo (0.1 sec lsb),
		7 - Split/Interval Time Mid,
		8 - Split/Interval Time High,
		9 - Split/Interval Distance Lo ( 1m lsb),
		10 - Split/Interval Distance Mid,
		11 - Split/Interval Distance High,
		12 - Interval Rest Time Lo (1 sec lsb),
		13 - Interval Rest Time Hi,
		14 - Interval Rest Distance Lo (1m lsb),
		15 - Interval Rest Distance Hi
		16 - Split/Interval Type10,
		17 - Split/Interval Number,	
*/
func (eng rowingEngine) GenerateSplitIntervalChar() []byte {
	bytes := make([]byte, 18, 18)

	copy(bytes[0:], eng.generateElapsedTimeBytes())
	return bytes
}

/*
	Generates a 'Split/Interval 2' characteristic of 19 bytes.
	Only the bytes labeled with a * are populated with non-zero values:
		0 - Elapsed Time Lo (0.01 sec lsb), *
		1 - Elapsed Time Mid,*
		2 - Elapsed Time High,*
		3 - Split/Interval Avg Stroke Rate,
		4 - Split/Interval Work Heartrate,
		5 - Split/Interval Rest Heartrate,
		6 - Split/Interval Avg Pace Lo (0.1 sec lsb)
		7 - Split/Interval Avg Pace Hi,
		8 - Split/Interval Total Calories Lo (Cals),
		9 - Split/Interval Total Calories Hi,
		10 - Split/Interval Avg Calories Lo (Cals/Hr),
		11 - Split/Interval Avg Calories Hi,
		12 - Split/Interval Speed Lo (0.001 m/s, max=65.534 m/s)
		13 - Split/Interval Speed Hi,
		14 - Split/Interval Power Lo (Watts, max = 65.534 kW)
		15 - Split/Interval Power Hi
		16 - Split Avg Drag Factor,
		17 - Split/Interval Number,
		18 - Erg Machine Type
*/
func (eng rowingEngine) GenerateSplitInterval2Char() []byte {
	bytes := make([]byte, 19, 19)

	copy(bytes[0:], eng.generateElapsedTimeBytes())
	return bytes
}

/*
	Generates a 'End of workout summary' characteristic of 20 bytes.
	Only the bytes labeled with a * are populated with non-zero values:
		0 - Log Entry Date Lo,
		1 - Log Entry Date Hi,
		2 - Log Entry Time Lo,
		3 - Log Entry Time Hi,
		4 - Elapsed Time Lo (0.01 sec lsb),
		5 - Elapsed Time Mid,
		6 - Elapsed Time High,
		7 - Distance Lo (0.1 m lsb),
		8 - Distance Mid,
		9 - Distance High,
		10 - Average Stroke Rate,
		11 - Ending Heartrate,
		12 - Average Heartrate,
		13 - Min Heartrate,
		14 - Max Heartrate,
		15 - Drag Factor Average,
		16 - Recovery Heart Rate, (zero = not valid data. After 1 minute of rest/recovery, PM5 resends this data unless the monitor has been turned off or a new workout started)
		17 - Workout Type,
		18 - Avg Pace Lo (0.1 sec lsb)
		19 - Avg Pace Hi
*/
func (eng rowingEngine) GenerateWorkoutSummaryChar() []byte {
	bytes := make([]byte, 20, 20)

	return bytes
}


/*
	Generates a 'End of workout summary 2' characteristic of 19 bytes.
	Only the bytes labeled with a * are populated with non-zero values:
		0 - Log Entry Date Lo,
		1 - Log Entry Date Hi,
		2 - Log Entry Time Lo,
		3 - Log Entry Time Hi,
		4 - Split/Interval Type12,
		5 - Split/Interval Size Lo, (meters or seconds)
		6 - Split/Interval Size Hi,
		7 - Split/Interval Count,
		8 - Total Calories Lo,
		9 - Total Calories Hi,
		10 - Watts Lo,
		11 - Watts Hi,
		12 - Total Rest Distance Lo (1 m lsb),
		13 - Total Rest Distance Mid,
		14 - Total Rest Distance High
		15 - Interval Rest Time Lo (seconds),
		16 - Interval Rest Time Hi,
		17 - Avg Calories Lo, (cals/hr)
		18 - Avg Calories Hi,
*/
func (eng rowingEngine) GenerateWorkoutSummary2Char() []byte {
	bytes := make([]byte, 19, 19)

	return bytes
}


/*
	Generates a 'Force curve data' characteristic of 2-288 bytes.
	Only the bytes labeled with a * are populated with non-zero values:
		0 - MS Nib = # characteristics, LS Nib = # words *
		1 - Sequence number, *
		2 - Data[n] (LS),
		3 - Data[n+1] (MS),
		4 - Data[n+2] (LS),
		5 - Data[n+3] (MS),
		6 - Data[n+4] (LS),
		7 - Data[n+5] (MS),
		8 - Data[n+6] (LS),
		9 - Data[n+7] (MS),
		10 - Data[n+8] (LS),
		11 - Data[n+9] (MS),
		12 - Data[n+10] (LS),
		13 - Data[n+11] (MS),
		14 - Data[n+12] (LS),
		15 - Data[n+13] (MS),
		16 - Data[n+14] (LS),
		17 - Data[n+15] (MS),
		18 - Data[n+16] (LS),
		19 - Data[n+17] (MS)
*/
func (eng rowingEngine) GenerateForceCurveChar() []byte {
	bytes := make([]byte, 20, 20)
	bytes[0:] = 0b00011001
	bytes[1:] = 1
	return bytes
}