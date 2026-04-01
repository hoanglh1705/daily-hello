import 'package:flutter_test/flutter_test.dart';
import 'package:daily_hello_app/core/utils/gps_fraud_detector.dart';

GpsFraudDetectionResult _makeResult({
  bool isMockLocation = false,
  bool isDeveloperModeOn = false,
  bool hasFakeGpsApps = false,
  List<String> detectedFakeApps = const [],
  bool sensorMismatch = false,
  bool isVpnActive = false,
  bool isDeviceCompromised = false,
}) {
  return GpsFraudDetectionResult(
    isMockLocation: isMockLocation,
    isDeveloperModeOn: isDeveloperModeOn,
    hasFakeGpsApps: hasFakeGpsApps,
    detectedFakeApps: detectedFakeApps,
    sensorMismatch: sensorMismatch,
    isVpnActive: isVpnActive,
    isDeviceCompromised: isDeviceCompromised,
  );
}

void main() {
  group('GpsFraudDetectionResult', () {
    test('isFraudulent returns true when mock location detected', () {
      final result = _makeResult(isMockLocation: true);
      expect(result.isFraudulent, isTrue);
      expect(result.reason, contains('Mock Location'));
    });

    test('isFraudulent returns true when fake GPS apps detected', () {
      final result = _makeResult(
        hasFakeGpsApps: true,
        detectedFakeApps: ['com.lexa.fakegps'],
      );
      expect(result.isFraudulent, isTrue);
      expect(result.reason, contains('com.lexa.fakegps'));
    });

    test('isFraudulent returns true when sensor mismatch detected', () {
      final result = _makeResult(sensorMismatch: true);
      expect(result.isFraudulent, isTrue);
      expect(result.reason, contains('Cảm biến'));
    });

    test('isFraudulent returns true when VPN is active', () {
      final result = _makeResult(isVpnActive: true);
      expect(result.isFraudulent, isTrue);
      expect(result.reason, contains('VPN/Proxy'));
    });

    test('isFraudulent returns true when device is compromised', () {
      final result = _makeResult(isDeviceCompromised: true);
      expect(result.isFraudulent, isTrue);
      // On test platform (not Android/iOS), the message still contains the keyword
      expect(result.reason, contains('không an toàn'));
    });

    test('isFraudulent returns false when all checks pass', () {
      final result = _makeResult();
      expect(result.isFraudulent, isFalse);
      expect(result.reason, isNull);
    });

    test('reason combines multiple fraud indicators', () {
      final result = _makeResult(
        isMockLocation: true,
        isDeveloperModeOn: true,
        hasFakeGpsApps: true,
        detectedFakeApps: ['com.lexa.fakegps', 'com.fakegps.mock'],
        sensorMismatch: true,
        isVpnActive: true,
        isDeviceCompromised: true,
      );
      expect(result.isFraudulent, isTrue);
      final reason = result.reason!;
      expect(reason, contains('Mock Location'));
      expect(reason, contains('com.lexa.fakegps'));
      expect(reason, contains('com.fakegps.mock'));
      expect(reason, contains('Cảm biến'));
      expect(reason, contains('VPN/Proxy'));
      expect(reason, contains('không an toàn'));
    });

    test('detectedFakeApps defaults to empty list', () {
      final result = _makeResult();
      expect(result.detectedFakeApps, isEmpty);
    });

    test('VPN alone blocks check-in', () {
      final result = _makeResult(isVpnActive: true);
      expect(result.isFraudulent, isTrue);
      expect(result.isMockLocation, isFalse);
      expect(result.hasFakeGpsApps, isFalse);
      expect(result.sensorMismatch, isFalse);
      expect(result.isDeviceCompromised, isFalse);
    });

    test('device compromised alone blocks check-in', () {
      final result = _makeResult(isDeviceCompromised: true);
      expect(result.isFraudulent, isTrue);
      expect(result.isMockLocation, isFalse);
      expect(result.hasFakeGpsApps, isFalse);
      expect(result.isVpnActive, isFalse);
    });
  });
}
