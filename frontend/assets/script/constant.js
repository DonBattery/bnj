"use strict";

// Client's unique ID
const ClientID = Math.random().toString(36).substring(2) + (new Date()).getTime().toString(36);

function getDefaultLevel() {
  let stringMap = "1110000000000000000000" + "1000000000001000011000" + "1000111100001100000000" + "1000000000011110000011" + "1100000000111000000001" + "1110001111110000000001" + "1000000000000011110001" + "1000000000000000000011" + "1110011100000000000111" + "1000000000003100000001" + "1000000000031110000001" + "1011110000311111111001" + "1000000000000000000001" + "1100000000000000000011" + "2222222214000001333111" + "1111111111111111111111";
  let intMap = [];
  for (let i = 0; i != 22 * 16; ++i) {
      intMap[i] = stringMap[i] * 1;
  };
  return intMap;
};

const Dott = [
  {"pos":[  0, 0],"dims":[13,15],"hotspot":[-2,-1]},
  {"pos":[ 19, 0],"dims":[14,14],"hotspot":[-2, 0]},
  {"pos":[ 38, 0],"dims":[14,16],"hotspot":[-2, 2]},
  {"pos":[ 57, 0],"dims":[13,16],"hotspot":[-3, 2]},
  {"pos":[ 76, 0],"dims":[14,14],"hotspot":[-3,-2]},
  {"pos":[ 95, 0],"dims":[15,14],"hotspot":[-1, 0]},
  {"pos":[114, 0],"dims":[16,16],"hotspot":[ 0, 0]},
  {"pos":[133, 0],"dims":[17,16],"hotspot":[ 1, 0]},
  {"pos":[152, 0],"dims":[17,13],"hotspot":[-1,-3]},
  {"pos":[171, 0],"dims":[13,15],"hotspot":[-1,-1]},
  {"pos":[190, 0],"dims":[14,14],"hotspot":[ 0, 0]},
  {"pos":[209, 0],"dims":[14,16],"hotspot":[ 0, 2]},
  {"pos":[228, 0],"dims":[13,16],"hotspot":[ 0, 2]},
  {"pos":[247, 0],"dims":[14,14],"hotspot":[ 1,-1]},
  {"pos":[266, 0],"dims":[15,14],"hotspot":[ 0, 1]},
  {"pos":[285, 0],"dims":[16,16],"hotspot":[ 0, 2]},
  {"pos":[304, 0],"dims":[17,16],"hotspot":[ 0, 2]},
  {"pos":[323, 0],"dims":[17,13],"hotspot":[ 2,-4]},
];

const Jiffy = [
  {"pos":[342, 0],"dims":[13,15],"hotspot":[-2,-1]},
  {"pos":[361, 0],"dims":[14,14],"hotspot":[-2, 0]},
  {"pos":[380, 0],"dims":[14,16],"hotspot":[-2, 2]},
  {"pos":[  0,18],"dims":[13,16],"hotspot":[-3, 2]},
  {"pos":[ 19,18],"dims":[14,14],"hotspot":[-3,-2]},
  {"pos":[ 38,18],"dims":[15,14],"hotspot":[-1, 0]},
  {"pos":[ 57,18],"dims":[16,16],"hotspot":[ 0, 0]},
  {"pos":[ 76,18],"dims":[17,16],"hotspot":[ 1, 0]},
  {"pos":[ 95,18],"dims":[17,13],"hotspot":[-1,-3]},
  {"pos":[114,18],"dims":[13,15],"hotspot":[-1,-1]},
  {"pos":[133,18],"dims":[14,14],"hotspot":[ 0, 0]},
  {"pos":[152,18],"dims":[14,16],"hotspot":[ 0, 2]},
  {"pos":[171,18],"dims":[13,16],"hotspot":[ 0, 2]},
  {"pos":[190,18],"dims":[14,14],"hotspot":[ 1,-1]},
  {"pos":[209,18],"dims":[15,14],"hotspot":[ 0, 1]},
  {"pos":[228,18],"dims":[16,16],"hotspot":[ 0, 2]},
  {"pos":[247,18],"dims":[17,16],"hotspot":[ 0, 2]},
  {"pos":[266,18],"dims":[17,13],"hotspot":[ 2,-4]},
];

const Fizz = [
  {"pos":[285,18],"dims":[13,15],"hotspot":[-2,-1]},
  {"pos":[304,18],"dims":[14,14],"hotspot":[-2, 0]},
  {"pos":[323,18],"dims":[14,16],"hotspot":[-2, 2]},
  {"pos":[342,18],"dims":[13,16],"hotspot":[-3, 2]},
  {"pos":[361,18],"dims":[14,14],"hotspot":[-3,-2]},
  {"pos":[380,18],"dims":[15,14],"hotspot":[-1, 0]},
  {"pos":[  0,36],"dims":[16,16],"hotspot":[ 0, 0]},
  {"pos":[ 19,36],"dims":[17,16],"hotspot":[ 1, 0]},
  {"pos":[ 38,36],"dims":[17,13],"hotspot":[-1,-3]},
  {"pos":[ 57,36],"dims":[13,15],"hotspot":[-1,-1]},
  {"pos":[ 76,36],"dims":[14,14],"hotspot":[ 0, 0]},
  {"pos":[ 95,36],"dims":[14,16],"hotspot":[ 0, 2]},
  {"pos":[114,36],"dims":[13,16],"hotspot":[ 0, 2]},
  {"pos":[133,36],"dims":[14,14],"hotspot":[ 1,-1]},
  {"pos":[152,36],"dims":[15,14],"hotspot":[ 0, 1]},
  {"pos":[171,36],"dims":[16,16],"hotspot":[ 0, 2]},
  {"pos":[190,36],"dims":[17,16],"hotspot":[ 0, 2]},
  {"pos":[209,36],"dims":[17,13],"hotspot":[ 2,-4]},
];

const Mijji = [
  {"pos":[228,36],"dims":[13,15],"hotspot":[-2,-1]},
  {"pos":[247,36],"dims":[14,14],"hotspot":[-2,-1]},
  {"pos":[266,36],"dims":[14,16],"hotspot":[-2, 2]},
  {"pos":[285,36],"dims":[13,16],"hotspot":[-3, 2]},
  {"pos":[304,36],"dims":[14,14],"hotspot":[-3,-2]},
  {"pos":[323,36],"dims":[15,14],"hotspot":[-1, 0]},
  {"pos":[342,36],"dims":[16,16],"hotspot":[ 0, 0]},
  {"pos":[361,36],"dims":[17,16],"hotspot":[ 1, 0]},
  {"pos":[380,36],"dims":[17,13],"hotspot":[-1,-3]},
  {"pos":[  0,54],"dims":[13,15],"hotspot":[-1,-1]},
  {"pos":[ 19,54],"dims":[14,14],"hotspot":[ 0, 0]},
  {"pos":[ 38,54],"dims":[14,16],"hotspot":[ 0, 2]},
  {"pos":[ 57,54],"dims":[13,16],"hotspot":[ 0, 2]},
  {"pos":[ 76,54],"dims":[14,14],"hotspot":[ 1,-1]},
  {"pos":[ 95,54],"dims":[15,14],"hotspot":[ 0, 1]},
  {"pos":[114,54],"dims":[16,16],"hotspot":[ 0, 2]},
  {"pos":[133,54],"dims":[17,16],"hotspot":[ 0, 2]},
  {"pos":[152,54],"dims":[17,13],"hotspot":[ 2,-4]}
];
