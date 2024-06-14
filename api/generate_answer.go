package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"golang.org/x/crypto/sha3"
)

var (
	cores   = []int{8, 12, 16, 24}
	screens = []int{3000, 4000, 6000}
)

var navigatorFp = []string{
	"vendorSub−",
	"productSub−20030107",
	"vendor−Google Inc.",
	"maxTouchPoints−0",
	"scheduling−[object Scheduling]",
	"userActivation−[object UserActivation]",
	"doNotTrack",
	"geolocation−[object Geolocation]",
	"connection−[object NetworkInformation]",
	"plugins−[object PluginArray]",
	"mimeTypes−[object MimeTypeArray]",
	"pdfViewerEnabled−true",
	"webkitTemporaryStorage−[object DeprecatedStorageQuota]",
	"webkitPersistentStorage−[object DeprecatedStorageQuota]",
	"windowControlsOverlay−[object WindowControlsOverlay]",
	"hardwareConcurrency−8",
	"cookieEnabled−true",
	"appCodeName−Mozilla",
	"appName−Netscape",
	//"appVersion−5.0 (Macintosh; Intel Mac OS X 10_15_7)…KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	"platform−MacIntel",
	"product−Gecko",
	//"userAgent−Mozilla/5.0 (Macintosh; Intel Mac OS X 1…KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	"language−zh-CN",
	"languages−zh-CN,zh",
	"onLine−true",
	"webdriver−false",
	"getGamepads−function getGamepads() { [native code] }",
	"javaEnabled−function javaEnabled() { [native code] }",
	"sendBeacon−function sendBeacon() { [native code] }",
	"vibrate−function vibrate() { [native code] }",
	"deprecatedRunAdAuctionEnforcesKAnonymity−false",
	"bluetooth−[object Bluetooth]",
	"storageBuckets−[object StorageBucketManager]",
	"clipboard−[object Clipboard]",
	"credentials−[object CredentialsContainer]",
	"keyboard−[object Keyboard]",
	"managed−[object NavigatorManagedData]",
	"mediaDevices−[object MediaDevices]",
	"storage−[object StorageManager]",
	"serviceWorker−[object ServiceWorkerContainer]",
	"virtualKeyboard−[object VirtualKeyboard]",
	"wakeLock−[object WakeLock]",
	"deviceMemory−8",
	"userAgentData−[object NavigatorUAData]",
	"login−[object NavigatorLogin]",
	"ink−[object Ink]",
	"mediaCapabilities−[object MediaCapabilities]",
	"hid−[object HID]",
	"locks−[object LockManager]",
	"gpu−[object GPU]",
	"mediaSession−[object MediaSession]",
	"permissions−[object Permissions]",
	"presentation−[object Presentation]",
	"usb−[object USB]",
	"xr−[object XRSystem]",
	"serial−[object Serial]",
	"adAuctionComponents−function adAuctionComponents() { [native code] }",
	"runAdAuction−function runAdAuction() { [native code] }",
	"canLoadAdAuctionFencedFrame−function canLoadAdAuctionFencedFrame() { [native code] }",
	"clearAppBadge−function clearAppBadge() { [native code] }",
	"getBattery−function getBattery() { [native code] }",
	"getUserMedia−function getUserMedia() { [native code] }",
	"requestMIDIAccess−function requestMIDIAccess() { [native code] }",
	"requestMediaKeySystemAccess−function requestMediaKeySystemAccess() { [native code] }",
	"setAppBadge−function setAppBadge() { [native code] }",
	"webkitGetUserMedia−function webkitGetUserMedia() { [native code] }",
	"clearOriginJoinedAdInterestGroups−function clearOriginJoinedAdInterestGroups() { [native code] }",
	"createAuctionNonce−function createAuctionNonce() { [native code] }",
	"deprecatedReplaceInURN−function deprecatedReplaceInURN() { [native code] }",
	"deprecatedURNToURL−function deprecatedURNToURL() { [native code] }",
	"getInstalledRelatedApps−function getInstalledRelatedApps() { [native code] }",
	"joinAdInterestGroup−function joinAdInterestGroup() { [native code] }",
	"leaveAdInterestGroup−function leaveAdInterestGroup() { [native code] }",
	"updateAdInterestGroups−function updateAdInterestGroups() { [native code] }",
	"registerProtocolHandler−function registerProtocolHandler() { [native code] }",
	"unregisterProtocolHandler−function unregisterProtocolHandler() { [native code] }",
}

var documentKeys = []string{
	"location",
	"_reactListeningebqqk1eqhpo",
}

var jsGlobals = []string{"0", "1", "window", "self", "document", "name", "location", "customElements", "history", "navigation", "locationbar", "menubar", "personalbar", "scrollbars", "statusbar", "toolbar", "status", "closed", "frames", "length", "top", "opener", "parent", "frameElement", "navigatorFp", "origin", "external", "screen", "innerWidth", "innerHeight", "scrollX", "pageXOffset", "scrollY", "pageYOffset", "visualViewport", "screenX", "screenY", "outerWidth", "outerHeight", "devicePixelRatio", "clientInformation", "screenLeft", "screenTop", "styleMedia", "onsearch", "isSecureContext", "trustedTypes", "performance", "onappinstalled", "onbeforeinstallprompt", "crypto", "indexedDB", "sessionStorage", "localStorage", "onbeforexrselect", "onabort", "onbeforeinput", "onbeforematch", "onbeforetoggle", "onblur", "oncancel", "oncanplay", "oncanplaythrough", "onchange", "onclick", "onclose", "oncontentvisibilityautostatechange", "oncontextlost", "oncontextmenu", "oncontextrestored", "oncuechange", "ondblclick", "ondrag", "ondragend", "ondragenter", "ondragleave", "ondragover", "ondragstart", "ondrop", "ondurationchange", "onemptied", "onended", "onerror", "onfocus", "onformdata", "oninput", "oninvalid", "onkeydown", "onkeypress", "onkeyup", "onload", "onloadeddata", "onloadedmetadata", "onloadstart", "onmousedown", "onmouseenter", "onmouseleave", "onmousemove", "onmouseout", "onmouseover", "onmouseup", "onmousewheel", "onpause", "onplay", "onplaying", "onprogress", "onratechange", "onreset", "onresize", "onscroll", "onsecuritypolicyviolation", "onseeked", "onseeking", "onselect", "onslotchange", "onstalled", "onsubmit", "onsuspend", "ontimeupdate", "ontoggle", "onvolumechange", "onwaiting", "onwebkitanimationend", "onwebkitanimationiteration", "onwebkitanimationstart", "onwebkittransitionend", "onwheel", "onauxclick", "ongotpointercapture", "onlostpointercapture", "onpointerdown", "onpointermove", "onpointerrawupdate", "onpointerup", "onpointercancel", "onpointerover", "onpointerout", "onpointerenter", "onpointerleave", "onselectstart", "onselectionchange", "onanimationend", "onanimationiteration", "onanimationstart", "ontransitionrun", "ontransitionstart", "ontransitionend", "ontransitioncancel", "onafterprint", "onbeforeprint", "onbeforeunload", "onhashchange", "onlanguagechange", "onmessage", "onmessageerror", "onoffline", "ononline", "onpagehide", "onpageshow", "onpopstate", "onrejectionhandled", "onstorage", "onunhandledrejection", "onunload", "crossOriginIsolated", "scheduler", "alert", "atob", "blur", "btoa", "cancelAnimationFrame", "cancelIdleCallback", "captureEvents", "clearInterval", "clearTimeout", "close", "confirm", "createImageBitmap", "fetch", "find", "focus", "getComputedStyle", "getSelection", "matchMedia", "moveBy", "moveTo", "open", "postMessage", "print", "prompt", "queueMicrotask", "releaseEvents", "reportError", "requestAnimationFrame", "requestIdleCallback", "resizeBy", "resizeTo", "scroll", "scrollBy", "scrollTo", "setInterval", "setTimeout", "stop", "structuredClone", "webkitCancelAnimationFrame", "webkitRequestAnimationFrame", "chrome", "fence", "caches", "cookieStore", "ondevicemotion", "ondeviceorientation", "ondeviceorientationabsolute", "launchQueue", "sharedStorage", "documentPictureInPicture", "getScreenDetails", "queryLocalFonts", "showDirectoryPicker", "showOpenFilePicker", "showSaveFilePicker", "originAgentCluster", "onpageswap", "onpagereveal", "credentialless", "speechSynthesis", "onscrollend", "webkitRequestFileSystem", "webkitResolveLocalFileSystemURL", "_sentryDebugIds", "webpackChunk_N_E", "__next_set_public_path__", "next", "__NEXT_DATA__", "__SSG_MANIFEST_CB", "__NEXT_P", "_N_E", "regeneratorRuntime", "__REACT_INTL_CONTEXT__", "DD_RUM", "_", "filterCSS", "filterXSS", "__SEGMENT_INSPECTOR__", "__NEXT_PRELOADREADY", "Intercom", "__MIDDLEWARE_MATCHERS", "__BUILD_MANIFEST", "__SSG_MANIFEST", "__STATSIG_SDK__", "__STATSIG_JS_SDK__", "__STATSIG_RERENDER_OVERRIDE__", "_oaiHandleSessionExpired", "__intercomAssignLocation", "__intercomReloadLocation"}

func BytesCombine2(pBytes ...[]byte) []byte {
	totalLen := 0
	for _, b := range pBytes {
		totalLen += len(b)
	}

	buffer := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buffer)
	buffer.Reset()
	buffer.Grow(totalLen)

	for _, b := range pBytes {
		buffer.Write(b)
	}

	return buffer.Bytes()
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func getParseTime() string {
	nowUTC := time.Now().UTC()

	nowChina := nowUTC.Add(8 * time.Hour)
	formattedTime := nowChina.Format("Mon Jan 02 2006 15:04:05 GMT+0800 (中国标准时间)")

	return formattedTime
}

// 定义一个程序启动时的时间
var startTime = time.Now()

func refreshStartTime() {
	startTime = time.Now()
}

// performanceNow 返回自程序启动以来的时间，单位为毫秒
func performanceNow() int64 {
	// 获取当前时间
	now := time.Now()

	if now.Sub(startTime).Minutes() > 30 {
		refreshStartTime()
	}

	return now.Sub(startTime).Milliseconds() + int64(rand.Intn(99999))
}

func generateAnswer(payload RequestData) string {
	seed, difficulty := payload.Seed, payload.Difficulty
	config := getConfig(payload)
	diffLen := len(difficulty)
	hasher := sha3.New512()
	startTime := time.Now()

	// 序列化固定config
	configStr, _ := sonic.Marshal(config)
	configStrTemplate := string(configStr)

	// 找出3和9的起止位置,逗号防止匹配错误
	format := ",%d,"
	pos3Start := strings.Index(configStrTemplate, fmt.Sprintf(format, config[3]))
	pos3End := pos3Start + len(fmt.Sprintf(format, config[3]))
	pos9Start := strings.Index(configStrTemplate, fmt.Sprintf(format, config[9]))
	pos9End := pos9Start + len(fmt.Sprintf(format, config[9]))

	// 将不变的部分拆分为固定的字符串片段
	beforePos3 := configStrTemplate[:pos3Start]
	betweenPos3AndPos9 := configStrTemplate[pos3End:pos9Start]
	afterPos9 := configStrTemplate[pos9End:]

	// 刚好每组三个字节, 可直接拼接; 如果有超出存入bufPart准备与动态插入的值进行拼接
	i := len([]byte(beforePos3)) % 3
	var beforePos3Base string
	var bufPart1 []byte
	if i == 0 {
		beforePos3Base = base64.StdEncoding.EncodeToString([]byte(beforePos3))
		bufPart1 = []byte("")
	} else {
		beforePos3Base = base64.StdEncoding.EncodeToString([]byte(beforePos3[:len([]byte(beforePos3))-i]))
		bufPart1 = []byte(beforePos3[len([]byte(beforePos3))-i:])
	}

	// 准备三种偏移版本, 应对 bufPart1+insert1 需要填充的情况, 填充值存储在 bufBeforePart2 索引为填充数量, bufAfterPart2为超出部分等待与insert2拼接
	var base64part2 []string = make([]string, 3)
	var bufBeforePart2 [][]byte = make([][]byte, 3)
	var bufAfterPart2 [][]byte = make([][]byte, 3)
	for j := 0; j < 3; j++ {
		i = len([]byte(betweenPos3AndPos9)[j:]) % 3
		if i == 0 {
			base64part2[j] = base64.StdEncoding.EncodeToString([]byte(betweenPos3AndPos9)[j:])
			bufBeforePart2[j] = []byte(betweenPos3AndPos9)[:j]
			bufAfterPart2[j] = []byte(betweenPos3AndPos9[len([]byte(betweenPos3AndPos9))-i:])
		} else {
			base64part2[j] = base64.StdEncoding.EncodeToString([]byte(betweenPos3AndPos9)[j : len([]byte(betweenPos3AndPos9))-i])
			bufBeforePart2[j] = []byte(betweenPos3AndPos9)[:j]
			bufAfterPart2[j] = []byte(betweenPos3AndPos9[len([]byte(betweenPos3AndPos9))-i:])
		}
	}

	var base64part3 []string = make([]string, 3)
	var bufBeforePart3 [][]byte = make([][]byte, 3)
	for j := 0; j < 3; j++ {
		base64part3[j] = base64.StdEncoding.EncodeToString([]byte(afterPos9)[j:])
		bufBeforePart3[j] = []byte(afterPos9)[:j]
	}
	// 预估所需的最大字符串长度
	initialSize := len(beforePos3Base) + len(betweenPos3AndPos9) + len(afterPos9) + 24

	var buffer bytes.Buffer
	buffer.Grow(initialSize)
	for i := 0; i < 500000; i++ {
		config[3] = i
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		config[9] = elapsed.Milliseconds()
		if elapsed >= time.Duration(10)*time.Second {
			break
		}
		buffer.Reset()
		buffer.WriteString(",")
		buffer.WriteString(strconv.Itoa(config[3].(int)))
		buffer.WriteString(",")
		config3Str := buffer.String()

		buffer.Reset()
		buffer.WriteString(",")
		buffer.WriteString(strconv.FormatInt(config[9].(int64), 10))
		buffer.WriteString(",")
		config9Str := buffer.String()
		insert1 := BytesCombine2(bufPart1, []byte(config3Str))
		var insertBase1 string
		var insertBase2 string
		var betweenPos3AndPos9Base string
		var afterPos9Base string
		var insert2 []byte
		i := len(insert1) % 3
		if i == 0 {
			insertBase1 = base64.StdEncoding.EncodeToString(insert1)
			betweenPos3AndPos9Base = base64part2[0]
			insert2 = BytesCombine2(bufAfterPart2[0], []byte(config9Str))
		} else {
			insertBase1 = base64.StdEncoding.EncodeToString(BytesCombine2(insert1, bufBeforePart2[3-i]))
			betweenPos3AndPos9Base = base64part2[3-i]
			insert2 = BytesCombine2(bufAfterPart2[3-i], []byte(config9Str))
		}

		i = len(insert2) % 3
		if i == 0 {
			insertBase2 = base64.StdEncoding.EncodeToString(insert2)
			afterPos9Base = base64part3[0]
		} else {
			insertBase2 = base64.StdEncoding.EncodeToString(BytesCombine2(insert2, bufBeforePart3[3-i]))
			afterPos9Base = base64part3[3-i]
		}
		buffer.Reset()
		buffer.WriteString(beforePos3Base)
		buffer.WriteString(insertBase1)
		buffer.WriteString(betweenPos3AndPos9Base)
		buffer.WriteString(insertBase2)
		buffer.WriteString(afterPos9Base)
		base := buffer.String()
		hasher.Write([]byte(seed + base))
		hash := hasher.Sum(nil)
		hasher.Reset()
		if hex.EncodeToString(hash[:diffLen])[:diffLen] <= difficulty {
			return "gAAAAAB" + string(base)
		}
	}
	return "gAAAAABwQ8Lk5FbGpA2NcR9dShT6gYjU7VxZ4D" + base64.StdEncoding.EncodeToString([]byte(`"`+seed+`"`))
}

func getConfig(payload RequestData) []interface{} {
	rand.NewSource(time.Now().UnixNano())
	core := cores[rand.Intn(len(cores))]
	screen := screens[rand.Intn(len(screens))]
	nfp := navigatorFp[rand.Intn(len(navigatorFp))]
	dk := documentKeys[rand.Intn(len(documentKeys))]
	wk := jsGlobals[rand.Intn(len(jsGlobals))]
	return []interface{}{
		core + screen,
		getParseTime(),
		int64(4294705152),
		3,
		payload.UserAgent,
		payload.Script,
		payload.CachedDpl,
		"zh-CN",
		"zh-CN,en,en-GB,en-US",
		9,
		nfp,
		dk,
		wk,
		performanceNow(),
		payload.Sid,
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	var requestData RequestData
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if requestData.Seed == "" || requestData.Difficulty == "" || requestData.UserAgent == "" || requestData.Script == "" || requestData.CachedDpl == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	answer := generateAnswer(requestData)
	w.Write([]byte(answer))
}

type RequestData struct {
	Seed       string `json:"seed"`
	Difficulty string `json:"difficulty"`
	UserAgent  string `json:"user_agent"`
	Script     string `json:"script_src"`
	CachedDpl  string `json:"dpl"`
	Sid        string `json:"sid"`
}
