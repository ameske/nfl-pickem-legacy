#!/bin/bash

nfl scrape scores
nfl grade
nfl generate results
mv *.html /opt/ameske/gonfl/templates/
