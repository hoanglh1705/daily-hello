import Flutter
import UIKit

@main
@objc class AppDelegate: FlutterAppDelegate {
  override func application(
    _ application: UIApplication,
    didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?
  ) -> Bool {
    let controller = window?.rootViewController as! FlutterViewController
    let fraudChannel = FlutterMethodChannel(
      name: "com.dailyhello/fraud_detection",
      binaryMessenger: controller.binaryMessenger
    )

    fraudChannel.setMethodCallHandler { (call, result) in
      switch call.method {
      case "isDeviceCompromised":
        result(self.isJailbroken())
      default:
        result(FlutterMethodNotImplemented)
      }
    }

    GeneratedPluginRegistrant.register(with: self)
    return super.application(application, didFinishLaunchingWithOptions: launchOptions)
  }

  private func isJailbroken() -> Bool {
    #if targetEnvironment(simulator)
      return false
    #else
      // Check for common jailbreak files
      let jailbreakPaths = [
        "/Applications/Cydia.app",
        "/Applications/Sileo.app",
        "/Applications/Zebra.app",
        "/Library/MobileSubstrate/MobileSubstrate.dylib",
        "/bin/bash",
        "/usr/sbin/sshd",
        "/etc/apt",
        "/usr/bin/ssh",
        "/private/var/lib/apt/",
        "/private/var/lib/cydia",
        "/private/var/stash",
        "/private/var/mobile/Library/SBSettings/Themes",
        "/var/lib/cydia",
        "/usr/libexec/cydia",
      ]

      for path in jailbreakPaths {
        if FileManager.default.fileExists(atPath: path) {
          return true
        }
      }

      // Check if app can write outside sandbox
      let testPath = "/private/jailbreak_test.txt"
      do {
        try "test".write(toFile: testPath, atomically: true, encoding: .utf8)
        try FileManager.default.removeItem(atPath: testPath)
        return true
      } catch {
        // Expected on non-jailbroken device
      }

      // Check if Cydia URL scheme is available
      if let cydiaUrl = URL(string: "cydia://package/com.example.package"),
         UIApplication.shared.canOpenURL(cydiaUrl) {
        return true
      }

      return false
    #endif
  }
}
