
for i in Parse0 ProcessPath0 ProcessPath1 Scan0 Set01 Set02 Set03 Set04 Set05 Set06 Set07 Set08 Set09 Set10 Set11 Set12 SetDefaults0 SetDefaults1 Validate02 Validate20 ; do
	go test -run "$i" >,r_$i
	if grep PASS ,r_$i >/dev/null 2>&1 ; then
		rm ,r_$i
	else
		echo FAIL $i
	fi
done


