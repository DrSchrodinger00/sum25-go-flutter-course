// ignore_for_file: unused_local_variable

import 'package:flutter/material.dart';
import '../services/preferences_service.dart';
import '../services/database_service.dart';
import '../services/secure_storage_service.dart';

class HomeScreen extends StatefulWidget {
  @override
  _HomeScreenState createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  bool _isLoading = false;
  String _statusMessage = '';

  @override
  void initState() {
    super.initState();
    PreferencesService.init();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('Lab04 Demo')),
      body: Padding(
        padding: EdgeInsets.all(16),
        child: Column(
          children: [
            ElevatedButton(
              onPressed: _testPrefs,
              child: Text('Test SharedPreferences'),
            ),
            ElevatedButton(
              onPressed: _testSQLite,
              child: Text('Test SQLite'),
            ),
            ElevatedButton(
              onPressed: _testSecureStorage,
              child: Text('Test SecureStorage'),
            ),
            if (_isLoading) CircularProgressIndicator(),
            if (_statusMessage.isNotEmpty) Padding(
              padding: EdgeInsets.only(top: 16),
              child: Text(_statusMessage),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _testPrefs() async {
    setState(() {
      _isLoading = true;
      _statusMessage = 'Testing SharedPreferences…';
    });
    try {
      await PreferencesService.setString('test_key', 'Hello prefs');
      final v = PreferencesService.getString('test_key');
      _statusMessage = 'SharedPreferences: \$v';
    } catch (e) {
      _statusMessage = 'SharedPreferences error: \$e';
    } finally {
      setState(() => _isLoading = false);
    }
  }

  Future<void> _testSQLite() async {
    setState(() {
      _isLoading = true;
      _statusMessage = 'Testing SQLite…';
    });
    try {
      final count = await DatabaseService.getUserCount();
      _statusMessage = 'SQLite users count: \$count';
    } catch (e) {
      _statusMessage = 'SQLite error: \$e';
    } finally {
      setState(() => _isLoading = false);
    }
  }

  Future<void> _testSecureStorage() async {
    setState(() {
      _isLoading = true;
      _statusMessage = 'Testing SecureStorage…';
    });
    try {
      await SecureStorageService.saveSecureData('secret', '42');
      final v = await SecureStorageService.getSecureData('secret');
      _statusMessage = 'SecureStorage: \$v';
    } catch (e) {
      _statusMessage = 'SecureStorage error: \$e';
    } finally {
      setState(() => _isLoading = false);
    }
  }
}
