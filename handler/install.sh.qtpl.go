// Code generated by qtc from "install.sh.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line handler/install.sh.qtpl:1
package handler

//line handler/install.sh.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line handler/install.sh.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line handler/install.sh.qtpl:1
func StreamShell(qw422016 *qt422016.Writer, r Result) {
//line handler/install.sh.qtpl:1
	qw422016.N().S(`
#!/bin/bash
if [ "$DEBUG" == "1" ]; then
	set -x
fi
TMP_DIR=$(mktemp -d -t worker-installer-XXXXXXXXXX)
function cleanup {
	rm -rf $TMP_DIR > /dev/null
}
function fail {
	cleanup
	msg=$1
	echo "============"
	echo "Error: $msg" 1>&2
	exit 1
}
function choose_asset {
	OS_ARCH=$1
	DISTRO=$2
	case "${OS_ARCH}_${DISTRO}" in`)
//line handler/install.sh.qtpl:20
	for _, n := range r.Assets {
//line handler/install.sh.qtpl:20
		qw422016.N().S(`
	"`)
//line handler/install.sh.qtpl:21
		qw422016.E().S(n.OS)
//line handler/install.sh.qtpl:21
		qw422016.N().S(`_`)
//line handler/install.sh.qtpl:21
		qw422016.E().S(n.Arch)
//line handler/install.sh.qtpl:21
		qw422016.N().S(`_`)
//line handler/install.sh.qtpl:21
		qw422016.E().S(n.Distro)
//line handler/install.sh.qtpl:21
		qw422016.N().S(`")
		echo "`)
//line handler/install.sh.qtpl:22
		qw422016.E().S(n.URL)
//line handler/install.sh.qtpl:22
		qw422016.N().S(` `)
//line handler/install.sh.qtpl:22
		qw422016.E().S(n.Type)
//line handler/install.sh.qtpl:22
		qw422016.N().S(`"
		return
		;;`)
//line handler/install.sh.qtpl:24
	}
//line handler/install.sh.qtpl:24
	qw422016.N().S(`
	*) exit 1;;
	esac
}
function install {
	#settings
	USER="`)
//line handler/install.sh.qtpl:30
	qw422016.E().S(r.User)
//line handler/install.sh.qtpl:30
	qw422016.N().S(`"
	IFS=',' read -r -a PROG_LIST <<< "`)
//line handler/install.sh.qtpl:31
	qw422016.E().S(r.Program)
//line handler/install.sh.qtpl:31
	qw422016.N().S(`"
	ASPROG="`)
//line handler/install.sh.qtpl:32
	if len(r.AsProgram) > 0 {
//line handler/install.sh.qtpl:32
		qw422016.N().S(` `)
//line handler/install.sh.qtpl:32
		qw422016.E().S(r.AsProgram)
//line handler/install.sh.qtpl:32
		qw422016.N().S(` `)
//line handler/install.sh.qtpl:32
	}
//line handler/install.sh.qtpl:32
	qw422016.N().S(`"
	MOVE="`)
//line handler/install.sh.qtpl:33
	qw422016.E().V(r.MoveToPath)
//line handler/install.sh.qtpl:33
	qw422016.N().S(`"
	RELEASE="`)
//line handler/install.sh.qtpl:34
	qw422016.E().S(r.Release)
//line handler/install.sh.qtpl:34
	qw422016.N().S(`" # `)
//line handler/install.sh.qtpl:34
	qw422016.E().S(r.ResolvedRelease)
//line handler/install.sh.qtpl:34
	qw422016.N().S(`
	INSECURE="`)
//line handler/install.sh.qtpl:35
	qw422016.E().V(r.Insecure)
//line handler/install.sh.qtpl:35
	qw422016.N().S(`"
	OUT_DIR="$(pwd)"
`)
//line handler/install.sh.qtpl:37
	if r.MoveToPath {
//line handler/install.sh.qtpl:37
		qw422016.N().S(`	if [ -d "$HOME/.local/bin" ]; then
		OUT_DIR="$HOME/.local/bin"
	elif [ -d "/opt/local/bin" ]; then
		OUT_DIR="/opt/local/bin"
	fi
`)
//line handler/install.sh.qtpl:43
	}
//line handler/install.sh.qtpl:43
	qw422016.N().S(`	GH="https://github.com"
	#bash check
	[ ! "$BASH_VERSION" ] && fail "Please use bash instead"
	[ ! -d $OUT_DIR ] && fail "output directory missing: $OUT_DIR"
	#dependency check, assume we are a standard POISX machine
	which find > /dev/null || fail "find not installed"
	which xargs > /dev/null || fail "xargs not installed"
	which sort > /dev/null || fail "sort not installed"
	which tail > /dev/null || fail "tail not installed"
	which cut > /dev/null || fail "cut not installed"
	which du > /dev/null || fail "du not installed"
	#choose an HTTP client
	GET=""
	if which curl > /dev/null; then
		GET="curl"
		if [[ $INSECURE = "true" ]]; then GET="$GET --insecure"; fi
		GET="$GET --fail -# -L"
	elif which wget > /dev/null; then
		GET="wget"
		if [[ $INSECURE = "true" ]]; then GET="$GET --no-check-certificate"; fi
		GET="$GET -qO-"
	else
		fail "neither wget/curl are installed"
	fi
	#debug HTTP
	if [ "$DEBUG" == "1" ]; then
		GET="$GET -v"
	fi
	#optional auth to install from private repos
	#NOTE: this also needs to be set on your instance of installer
	AUTH="${GITHUB_TOKEN}"
	if [ ! -z "$AUTH" ]; then
		GET="$GET -H 'Authorization: $AUTH'"
	fi
	#find OS #TODO BSDs and other posixs
	case `)
//line handler/install.sh.qtpl:43
	qw422016.N().S("`")
//line handler/install.sh.qtpl:43
	qw422016.N().S(`uname -s`)
//line handler/install.sh.qtpl:43
	qw422016.N().S("`")
//line handler/install.sh.qtpl:43
	qw422016.N().S(` in
	Darwin) OS="darwin";;
	Linux) OS="linux";;
	*) fail "unknown os: $(uname -s)";;
	esac
	#find ARCH
	if uname -m | grep -E '(arm|arch)64' > /dev/null; then
		ARCH="arm64"
		`)
//line handler/install.sh.qtpl:87
	if !r.M1Asset {
//line handler/install.sh.qtpl:87
		qw422016.N().S(`
		# no m1 assets. if on mac arm64, rosetta allows fallback to amd64
		if [[ $OS = "darwin" ]]; then
			ARCH="amd64"
		fi
		`)
//line handler/install.sh.qtpl:92
	}
//line handler/install.sh.qtpl:92
	qw422016.N().S(`
	elif uname -m | grep 64 > /dev/null; then
		ARCH="amd64"
	elif uname -m | grep arm > /dev/null; then
		ARCH="arm" #TODO armv6/v7
	elif uname -m | grep 386 > /dev/null; then
		ARCH="386"
	else
		fail "unknown arch: $(uname -m)"
	fi
	#find Distro
	if [ -f /etc/os-release ]; then
		DISTRO=$(grep ^ID_LIKE= /etc/os-release | cut -d '=' -f 2-)
		if [ -z "$DISTRO" ]; then
			DISTRO=$(grep ^ID= /etc/os-release | cut -d '=' -f 2-)
		fi
	else
		DISTRO="generic"
	fi
	#choose from asset list
	OS_ARCH="${OS}_${ARCH}"
	ASSET_INFO=$(choose_asset "$OS_ARCH" "generic")
	if [ $? -ne 0 ]; then
		ASSET_INFO=$(choose_asset "$OS_ARCH" "$DISTRO")
	else
		DISTRO="generic"
	fi
	URL=$(echo $ASSET_INFO | cut -d ' ' -f 1)
	FTYPE=$(echo $ASSET_INFO | cut -d ' ' -f 2)
	if [ -z "$URL" ] || [ -z "$FTYPE" ]; then
		fail "No valid asset found for ${OS_ARCH}"
	fi
	#got URL! download it...
	echo -n "`)
//line handler/install.sh.qtpl:125
	if r.MoveToPath {
//line handler/install.sh.qtpl:125
		qw422016.N().S(`Installing`)
//line handler/install.sh.qtpl:125
	} else {
//line handler/install.sh.qtpl:125
		qw422016.N().S(`Downloading`)
//line handler/install.sh.qtpl:125
	}
//line handler/install.sh.qtpl:125
	qw422016.N().S(`"
	echo -n " $USER/${PROG_LIST[*]}"
	if [ ! -z "$RELEASE" ]; then
		echo -n " $RELEASE"
	fi
	if [ ! -z "$ASPROG" ]; then
		echo -n " as $ASPROG"
	fi
	echo -n " (${OS}/${ARCH})"
	`)
//line handler/install.sh.qtpl:134
	if r.Search {
//line handler/install.sh.qtpl:134
		qw422016.N().S(`
	# web search, give time to cancel
	echo -n " in 5 seconds"
	for i in 1 2 3 4 5; do
		sleep 1
		echo -n "."
	done
	`)
//line handler/install.sh.qtpl:141
	} else {
//line handler/install.sh.qtpl:141
		qw422016.N().S(`
	echo "....."
	`)
//line handler/install.sh.qtpl:143
	}
//line handler/install.sh.qtpl:143
	qw422016.N().S(`
	#enter tempdir
	mkdir -p $TMP_DIR
	cd $TMP_DIR
	if [[ $FTYPE = ".gz" ]]; then
		which gzip > /dev/null || fail "gzip is not installed"
		bash -c "$GET $URL" | gzip -d - > "${PROG_LIST[0]}" || fail "download failed"
	elif [[ $FTYPE = ".bz2" ]]; then
		which bzip2 > /dev/null || fail "bzip2 is not installed"
		bash -c "$GET $URL" | bzip2 -d - > "${PROG_LIST[0]}" || fail "download failed"
	elif [[ $FTYPE = ".tar.bz" ]] || [[ $FTYPE = ".tar.bz2" ]]; then
		which tar > /dev/null || fail "tar is not installed"
		which bzip2 > /dev/null || fail "bzip2 is not installed"
		bash -c "$GET $URL" | tar jxf - || fail "download failed"
	elif [[ $FTYPE = ".tar.gz" ]] || [[ $FTYPE = ".tgz" ]]; then
		which tar > /dev/null || fail "tar is not installed"
		which gzip > /dev/null || fail "gzip is not installed"
		bash -c "$GET $URL" | tar zxf - || fail "download failed"
	elif [[ $FTYPE = ".tar.xz" ]] || [[ $FTYPE = ".txz" ]]; then
		which tar > /dev/null || fail "tar is not installed"
		which xz > /dev/null || fail "xz is not installed"
		bash -c "$GET $URL" | tar Jxf - || fail "download failed"
	elif [[ $FTYPE = ".zip" ]]; then
		which unzip > /dev/null || fail "unzip is not installed"
		bash -c "$GET $URL" > tmp.zip || fail "download failed"
		unzip -o -qq tmp.zip || fail "unzip failed"
		rm tmp.zip || fail "cleanup failed"
	elif [[ $FTYPE = ".bin" ]]; then
		bash -c "$GET $URL" > "${PROG_LIST[0]}_${OS}_${ARCH}" || fail "download failed"
	elif [[ $FTYPE = ".deb" ]]; then
		which dpkg > /dev/null || fail "dpkg is not installed"
		bash -c "$GET $URL" > tmp.deb || fail "download failed"
		sudo dpkg -i tmp.deb || fail "dpkg install failed"
		rm tmp.deb || fail "cleanup failed"
	elif [[ $FTYPE = ".rpm" ]]; then
		which rpm > /dev/null || fail "rpm is not installed"
		bash -c "$GET $URL" > tmp.rpm || fail "download failed"
		sudo rpm -i tmp.rpm || fail "rpm install failed"
		rm tmp.rpm || fail "cleanup failed"
	else
		fail "unknown file type: $FTYPE"
	fi
	if [[ $DISTRO = "generic" ]]; then
		for PROG in "${PROG_LIST[@]}"; do
			BIN_PATH=$(find . -type f | grep -i "$PROG" | head -n 1)
			[[ -z "$BIN_PATH" ]] && fail "Binary $PROG not found"

			chmod +x "$BIN_PATH" || fail "chmod +x failed on $BIN_PATH"
			DEST="$OUT_DIR/$PROG"

			OUT=$(mv "$BIN_PATH" "$DEST" 2>&1)
			STATUS=$?
			if [ $STATUS -ne 0 ]; then
				if [[ $OUT =~ "Permission denied" ]]; then
					if [ -w "$DEST" ]; then
						mv "$BIN_PATH" "$DEST" || fail "mv failed for $BIN_PATH"
					else
						echo "mv with sudo..."
						sudo mv "$BIN_PATH" "$DEST" || fail "sudo mv failed for $BIN_PATH"
					fi
				else
					fail "mv failed for $BIN_PATH ($OUT)"
				fi
			fi
			echo "Moved $PROG to $DEST"
		done
	fi
	cleanup
}
install
`)
//line handler/install.sh.qtpl:213
}

//line handler/install.sh.qtpl:213
func WriteShell(qq422016 qtio422016.Writer, r Result) {
//line handler/install.sh.qtpl:213
	qw422016 := qt422016.AcquireWriter(qq422016)
//line handler/install.sh.qtpl:213
	StreamShell(qw422016, r)
//line handler/install.sh.qtpl:213
	qt422016.ReleaseWriter(qw422016)
//line handler/install.sh.qtpl:213
}

//line handler/install.sh.qtpl:213
func Shell(r Result) string {
//line handler/install.sh.qtpl:213
	qb422016 := qt422016.AcquireByteBuffer()
//line handler/install.sh.qtpl:213
	WriteShell(qb422016, r)
//line handler/install.sh.qtpl:213
	qs422016 := string(qb422016.B)
//line handler/install.sh.qtpl:213
	qt422016.ReleaseByteBuffer(qb422016)
//line handler/install.sh.qtpl:213
	return qs422016
//line handler/install.sh.qtpl:213
}
