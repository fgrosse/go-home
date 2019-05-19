.PHONY: setup

# The setup-fedora task installs all required non-go development dependencies
# This list is based on https://engoengine.github.io/tutorials/00-foreword which
# lists the dependencies for Ubuntu.
# Also https://github.com/go-gl/glfw/blob/master/README.md lists concrete
# package names for Fedora (CentOS)
setup-fedora:
	dnf install \
		alsa-lib-devel \
		mesa-libGLU-devel \
		mesa-libGL-devel \
		freeglut-devel \
		git-all \
		libX11-devel \
		libXcursor-devel \
		libXrandr-devel \
		libXinerama-devel \
		libXi-devel
