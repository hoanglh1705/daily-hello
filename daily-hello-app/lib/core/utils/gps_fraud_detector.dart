import 'dart:async';
import 'dart:io';
import 'dart:math';

import 'package:flutter/services.dart';
import 'package:geolocator/geolocator.dart';
import 'package:sensors_plus/sensors_plus.dart';

class GpsFraudDetectionResult {
  final bool isMockLocation;
  final bool isDeveloperModeOn;
  final bool hasFakeGpsApps;
  final List<String> detectedFakeApps;
  final bool sensorMismatch;
  final bool isVpnActive;
  final bool isDeviceCompromised;

  GpsFraudDetectionResult({
    required this.isMockLocation,
    required this.isDeveloperModeOn,
    required this.hasFakeGpsApps,
    this.detectedFakeApps = const [],
    required this.sensorMismatch,
    required this.isVpnActive,
    required this.isDeviceCompromised,
  });

  bool get isFraudulent =>
      isMockLocation ||
      hasFakeGpsApps ||
      sensorMismatch ||
      isVpnActive ||
      isDeviceCompromised;

  String? get reason {
    final reasons = <String>[];
    if (isMockLocation) reasons.add('Phát hiện vị trí giả lập (Mock Location)');
    if (hasFakeGpsApps) {
      reasons.add(
        'Phát hiện ứng dụng giả GPS: ${detectedFakeApps.join(", ")}',
      );
    }
    if (sensorMismatch) {
      reasons.add('Cảm biến thiết bị không khớp với vị trí GPS');
    }
    if (isVpnActive) {
      reasons.add('Phát hiện VPN/Proxy đang hoạt động');
    }
    if (isDeviceCompromised) {
      reasons.add(
        Platform.isAndroid
            ? 'Thiết bị đã bị root - không an toàn'
            : 'Thiết bị đã bị jailbreak - không an toàn',
      );
    }
    return reasons.isEmpty ? null : reasons.join('. ');
  }
}

class GpsFraudDetector {
  static const _channel = MethodChannel('com.dailyhello/fraud_detection');

  // Known fake GPS app package names
  static const _knownFakeGpsPackages = [
    'com.lexa.fakegps',
    'com.incorporateapps.fakegps.fre',
    'com.fakegps.mock',
    'com.evezzon.fakegps',
    'com.gsmartstudio.fakegps',
    'com.lkr.fakelocation',
    'com.marlon.floating.fake.location',
    'com.location.faker',
    'com.theappninjas.fakegpsjoystick',
    'com.incorporateapps.fakegps',
    'org.hola.gpslocation',
    'com.divi.fakeGPS',
    'fr.dvilleneuve.lockito',
    'com.rosteam.gpsemulator',
    'com.blogspot.newlooper.flp',
    'ru.gavrikov.mocklocations',
    'com.adjustyourgps.fake_gps_location',
  ];

  /// Run all fraud detection checks.
  /// [position] is the GPS position obtained from Geolocator.
  /// [sensorDuration] is how long to sample accelerometer data.
  static Future<GpsFraudDetectionResult> detect(
    Position position, {
    Duration sensorDuration = const Duration(seconds: 2),
  }) async {
    // Run checks in parallel
    final results = await Future.wait([
      _checkMockLocation(position),
      _checkDeveloperMode(),
      _checkFakeGpsApps(),
      _checkSensorConsistency(position, sensorDuration),
      _checkVpnActive(),
      _checkDeviceCompromised(),
    ]);

    return GpsFraudDetectionResult(
      isMockLocation: results[0] as bool,
      isDeveloperModeOn: results[1] as bool,
      hasFakeGpsApps: (results[2] as List<String>).isNotEmpty,
      detectedFakeApps: results[2] as List<String>,
      sensorMismatch: results[3] as bool,
      isVpnActive: results[4] as bool,
      isDeviceCompromised: results[5] as bool,
    );
  }

  /// Check if the position is from a mock provider.
  static Future<bool> _checkMockLocation(Position position) async {
    return position.isMocked;
  }

  /// Check if developer mode is enabled (Android only).
  static Future<bool> _checkDeveloperMode() async {
    if (!Platform.isAndroid) return false;
    try {
      final result = await _channel.invokeMethod<bool>('isDeveloperModeOn');
      return result ?? false;
    } on MissingPluginException {
      return false;
    }
  }

  /// Scan for known fake GPS apps (Android only).
  static Future<List<String>> _checkFakeGpsApps() async {
    if (!Platform.isAndroid) return [];
    try {
      final result = await _channel.invokeMethod<List<dynamic>>(
        'detectFakeGpsApps',
        {'packages': _knownFakeGpsPackages},
      );
      return result?.cast<String>() ?? [];
    } on MissingPluginException {
      return [];
    }
  }

  /// Check if sensor data is consistent with GPS position.
  ///
  /// If GPS reports near-zero speed (standing still) but the accelerometer
  /// shows significant movement beyond gravity, flag as suspicious.
  static Future<bool> _checkSensorConsistency(
    Position position,
    Duration duration,
  ) async {
    try {
      final samples = <_AccSample>[];
      final subscription = accelerometerEventStream(
        samplingPeriod: const Duration(milliseconds: 100),
      ).listen((event) {
        samples.add(_AccSample(event.x, event.y, event.z));
      });

      // Collect samples for the specified duration
      await Future.delayed(duration);
      await subscription.cancel();

      if (samples.length < 5) return false; // Not enough data

      // Calculate the magnitude of acceleration for each sample
      // Subtract gravity (~9.8) to get user acceleration
      final magnitudes = samples.map((s) {
        final magnitude = sqrt(s.x * s.x + s.y * s.y + s.z * s.z);
        return (magnitude - 9.81).abs();
      }).toList();

      // Calculate standard deviation of acceleration magnitude
      final mean = magnitudes.reduce((a, b) => a + b) / magnitudes.length;
      final variance = magnitudes
              .map((m) => (m - mean) * (m - mean))
              .reduce((a, b) => a + b) /
          magnitudes.length;
      final stdDev = sqrt(variance);

      // GPS says stationary (speed < 0.5 m/s) but accelerometer shows
      // significant variation (stdDev > 1.5), which indicates the device
      // is being physically moved while GPS is spoofed to stay in place.
      final gpsStationary = position.speed < 0.5;
      final sensorMoving = stdDev > 1.5;

      return gpsStationary && sensorMoving;
    } catch (_) {
      // If sensors are unavailable, skip this check
      return false;
    }
  }

  /// Detect active VPN/Proxy connections by checking network interfaces.
  /// VPN connections typically create tun0, ppp0, or similar interfaces.
  static Future<bool> _checkVpnActive() async {
    try {
      final interfaces = await NetworkInterface.list(
        includeLoopback: false,
        type: InternetAddressType.any,
      );

      // VPN interface name patterns
      const vpnPatterns = ['tun', 'tap', 'ppp', 'pptp', 'l2tp', 'ipsec', 'utun'];

      for (final iface in interfaces) {
        final name = iface.name.toLowerCase();
        for (final pattern in vpnPatterns) {
          if (name.startsWith(pattern)) {
            return true;
          }
        }
      }

      return false;
    } catch (_) {
      return false;
    }
  }

  /// Detect if the device is rooted (Android) or jailbroken (iOS).
  static Future<bool> _checkDeviceCompromised() async {
    try {
      final result =
          await _channel.invokeMethod<bool>('isDeviceCompromised');
      return result ?? false;
    } on MissingPluginException {
      return false;
    }
  }
}

class _AccSample {
  final double x, y, z;
  _AccSample(this.x, this.y, this.z);
}
