RELEASE_TARGETS := build_release build_strip build_upx build_sed build_info 
DEV_TARGETS = build_dev build_run build_remove
VALGRIND_TARGETS = build_dev build_valgrind

NINJA = ninja

.DEFAULT_GOAL = dev

release: 
	@$(foreach targ, $(RELEASE_TARGETS), \
		$(NINJA) $(targ); \
	)
dev: 
	@$(foreach targ, $(DEV_TARGETS), \
		$(NINJA) $(targ); \
	)
valgrind: 
	@$(foreach targ, $(VALGRIND_TARGETS), \
		$(NINJA) $(targ); \
	)	
cp: 
	docker run -it --name con cross 
	docker cp con:/home/user/gay.txt gay.txt
docker: 
	ninja -f cross-compile.ninja docker
run: 
	docker run -it --name con cross	
#
.PHONY: docker
