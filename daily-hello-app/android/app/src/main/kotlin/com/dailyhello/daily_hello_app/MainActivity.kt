package com.dailyhello.daily_hello_app

import android.content.pm.PackageManager
import android.os.Build
import android.provider.Settings
import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.MethodChannel
import java.io.File

class MainActivity : FlutterActivity() {
    private val CHANNEL = "com.dailyhello/fraud_detection"

    override fun configureFlutterEngine(flutterEngine: FlutterEngine) {
        super.configureFlutterEngine(flutterEngine)

        MethodChannel(flutterEngine.dartExecutor.binaryMessenger, CHANNEL)
            .setMethodCallHandler { call, result ->
                when (call.method) {
                    "isDeveloperModeOn" -> {
                        result.success(isDeveloperModeOn())
                    }
                    "detectFakeGpsApps" -> {
                        val packages = call.argument<List<String>>("packages") ?: emptyList()
                        result.success(detectFakeGpsApps(packages))
                    }
                    "isDeviceCompromised" -> {
                        result.success(isDeviceRooted())
                    }
                    else -> result.notImplemented()
                }
            }
    }

    private fun isDeveloperModeOn(): Boolean {
        return try {
            Settings.Global.getInt(
                contentResolver,
                Settings.Global.DEVELOPMENT_SETTINGS_ENABLED,
                0
            ) != 0
        } catch (e: Exception) {
            false
        }
    }

    private fun detectFakeGpsApps(knownPackages: List<String>): List<String> {
        val detected = mutableListOf<String>()
        val pm = packageManager

        for (pkg in knownPackages) {
            try {
                pm.getPackageInfo(pkg, PackageManager.GET_META_DATA)
                detected.add(pkg)
            } catch (e: PackageManager.NameNotFoundException) {
                // App not installed — expected
            }
        }

        return detected
    }

    private fun isDeviceRooted(): Boolean {
        return checkRootBinaries() || checkSuExists() || checkDangerousProps() || checkRootManagementApps()
    }

    private fun checkRootBinaries(): Boolean {
        val paths = arrayOf(
            "/system/bin/su",
            "/system/xbin/su",
            "/sbin/su",
            "/system/su",
            "/system/bin/.ext/.su",
            "/system/usr/we-need-root/su-backup",
            "/system/xbin/mu",
            "/data/local/xbin/su",
            "/data/local/bin/su",
            "/data/local/su",
        )
        return paths.any { File(it).exists() }
    }

    private fun checkSuExists(): Boolean {
        return try {
            val process = Runtime.getRuntime().exec(arrayOf("/system/xbin/which", "su"))
            val result = process.inputStream.bufferedReader().readText().trim()
            result.isNotEmpty()
        } catch (e: Exception) {
            false
        }
    }

    private fun checkDangerousProps(): Boolean {
        return try {
            val process = Runtime.getRuntime().exec(arrayOf("getprop", "ro.debuggable"))
            val result = process.inputStream.bufferedReader().readText().trim()
            result == "1"
        } catch (e: Exception) {
            false
        }
    }

    private fun checkRootManagementApps(): Boolean {
        val rootApps = arrayOf(
            "com.topjohnwu.magisk",
            "eu.chainfire.supersu",
            "com.koushikdutta.superuser",
            "com.noshufou.android.su",
            "com.thirdparty.superuser",
            "com.yellowes.su",
            "com.devadvance.rootcloak",
            "com.devadvance.rootcloakplus",
            "de.robv.android.xposed.installer",
            "com.saurik.substrate",
        )
        val pm = packageManager
        return rootApps.any { pkg ->
            try {
                pm.getPackageInfo(pkg, PackageManager.GET_META_DATA)
                true
            } catch (e: PackageManager.NameNotFoundException) {
                false
            }
        }
    }
}
