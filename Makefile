movies:
	@for name in [0-9][0-9][0-9][0-9]; do\
		ffmpeg -loglevel quiet -framerate 25 -i $${name}/image%03d.png -c:v libx264 -profile:v high -crf 20 -pix_fmt yuv420p $${name}.mp4;\
		echo $${name};\
		cp $${name}/peeps.txt $${name}.txt;\
		rm -rf $${name};\
	done

run:
	@go run ./src