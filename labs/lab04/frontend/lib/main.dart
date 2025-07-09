// lib/main.dart

import 'package:flutter/material.dart';
import 'package:lab04_frontend/services/preferences_service.dart';
import 'package:lab04_frontend/services/database_service.dart';
import 'package:lab04_frontend/screens/home_screen.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Initialize all persistence services
  await PreferencesService.init();
  await DatabaseService.database;

  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Lab 04 - Database & Persistence',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
        useMaterial3: true,
      ),
      home: HomeScreen(),
      debugShowCheckedModeBanner: false,
    );
  }
}
