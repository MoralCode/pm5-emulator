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
