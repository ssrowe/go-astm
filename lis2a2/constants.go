package lis2a2

type Encoding int

const EncodingUTF8 Encoding = 1
const EncodingASCII Encoding = 2
const EncodingWindows1250 Encoding = 3
const EncodingWindows1251 Encoding = 4
const EncodingWindows1252 Encoding = 5
const EncodingDOS852 Encoding = 6
const EncodingDOS855 Encoding = 7
const EncodingDOS866 Encoding = 8

type Timezone string

const TimezoneUTC Timezone = "UTC"
const TimezoneEuropeBerlin Timezone = "Europe/Berlin"
const TimezoneEuropeBudapest Timezone = "Europe/Budapest"
const TimezoneEuropeLondon Timezone = "Europe/London"

type LineBreak int

const CR LineBreak = 0x13
const LF LineBreak = 0x10
const CRLF LineBreak = 0x1310

type Notation int

const StandardNotation = 1
const ShortNotation = 2
