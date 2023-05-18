#!/bin/bash

# List all Go files excluding cmd, vendor, and pkg directories
list_files=$(go list ./... | grep -v /cmd/ | grep -v /vendor/ | grep -v mock | grep -v /scripts/)

# Run the Go tests and capture the output to a file
echo "=== RUNNING TEST ==="
go test -cover $list_files | tee test.out

# Extract the coverage statistics and sum them
success_count=0
coverage_sum=0
while read -r line; do
  if [[ $line =~ ok ]]; then
    success_count=$((success_count + 1))
    coverage=$(echo $line | awk '{print $5}')
    # remove percentage
    coverage_no_percent=${coverage//\%/}
    # split by . and got firest value
    coverage_no_percent=$(echo "$coverage_no_percent" | awk -F. '{ print $1 }')
    coverage_sum=$((coverage_sum+coverage_no_percent))
  elif [[ $line =~ FAILED ]]; then
    package=$(echo $line | awk '{print $2}')
    success="FAILED"
  elif [[ $line =~ \? ]]; then
    package=$(echo $line | awk '{print $2}')
    success="NO TEST FILES"
  fi
done < test.out

# Calculate the average coverage
list_count=$(echo "$list_files" | wc -l)
average_coverage=$((coverage_sum / list_count))

# Print the total coverage
echo ""
echo "=== TEST SUMMARY ==="
echo "Success: $success_count/$list_count"
echo "Average coverage: $average_coverage%"