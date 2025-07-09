import 'dart:convert';

class SecureStorageService {
  // In-memory storage for tests and runtime
  static final Map<String, String> _memory = {};

  static const _keyAuthToken = 'auth_token';
  static const _keyUsername = 'username';
  static const _keyPassword = 'password';
  static const _keyBiometric = 'biometric';

  static Future<void> saveAuthToken(String token) async {
    _memory[_keyAuthToken] = token;
  }

  static Future<String?> getAuthToken() async {
    return _memory[_keyAuthToken];
  }

  static Future<void> deleteAuthToken() async {
    _memory.remove(_keyAuthToken);
  }

  static Future<void> saveUserCredentials(String username, String password) async {
    _memory[_keyUsername] = username;
    _memory[_keyPassword] = password;
  }

  static Future<Map<String, String?>> getUserCredentials() async {
    return {
      'username': _memory[_keyUsername],
      'password': _memory[_keyPassword],
    };
  }

  static Future<void> deleteUserCredentials() async {
    _memory.remove(_keyUsername);
    _memory.remove(_keyPassword);
  }

  static Future<void> saveBiometricEnabled(bool enabled) async {
    _memory[_keyBiometric] = enabled.toString();
  }

  static Future<bool> isBiometricEnabled() async {
    return _memory[_keyBiometric] == 'true';
  }

  static Future<void> saveSecureData(String key, String value) async {
    _memory[key] = value;
  }

  static Future<String?> getSecureData(String key) async {
    return _memory[key];
  }

  static Future<void> deleteSecureData(String key) async {
    _memory.remove(key);
  }

  static Future<void> saveObject(String key, Map<String, dynamic> object) async {
    _memory[key] = jsonEncode(object);
  }

  static Future<Map<String, dynamic>?> getObject(String key) async {
    final jsonStr = _memory[key];
    if (jsonStr == null) return null;
    return jsonDecode(jsonStr) as Map<String, dynamic>;
  }

  static Future<bool> containsKey(String key) async {
    return _memory.containsKey(key);
  }

  static Future<List<String>> getAllKeys() async {
    return _memory.keys.toList();
  }

  static Future<void> clearAll() async {
    _memory.clear();
  }

  static Future<Map<String, String>> exportData() async {
    return Map<String, String>.from(_memory);
  }
}
