// lib/services/database_service.dart

import 'package:sqflite/sqflite.dart';
import 'package:sqflite/utils/utils.dart';
import 'package:sqflite_common_ffi/sqflite_ffi.dart';
import 'package:path/path.dart';
import '../models/user.dart';
import '../models/user.dart' show CreateUserRequest;

class DatabaseService {
  static Database? _database;
  static const String _dbName = 'lab04_app.db';
  static const int _version = 1;

  /// Returns the open database instance, initializing via ffi in tests.
  static Future<Database> get database async {
    if (_database != null) return _database!;
    sqfliteFfiInit();
    databaseFactory = databaseFactoryFfi;
    _database = await _initDatabase();
    return _database!;
  }

  static Future<Database> _initDatabase() async {
    final path = join(await getDatabasesPath(), _dbName);
    return openDatabase(
      path,
      version: _version,
      onCreate: _onCreate,
      onUpgrade: _onUpgrade,
    );
  }

  static Future<void> _onCreate(Database db, int version) async {
    await db.execute('''
      CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        email TEXT NOT NULL,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL
      );
    ''');
    await db.execute('''
      CREATE TABLE posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        title TEXT NOT NULL,
        content TEXT,
        published INTEGER NOT NULL,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL,
        FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
      );
    ''');
  }

  static Future<void> _onUpgrade(Database db, int oldVersion, int newVersion) async {
    // No migrations yet.
  }

  /// Inserts a new user and returns it with an assigned `id`.
  static Future<User> createUser(CreateUserRequest req) async {
    final db = await database;
    final now = DateTime.now().toIso8601String();
    final id = await db.insert('users', {
      'name': req.name,
      'email': req.email,
      'created_at': now,
      'updated_at': now,
    });
    return User(
      id: id,
      name: req.name,
      email: req.email,
      createdAt: DateTime.parse(now),
      updatedAt: DateTime.parse(now),
    );
  }

  /// Fetches a user by `id`, or returns `null` if none exists.
  static Future<User?> getUser(int id) async {
    final db = await database;
    final rows = await db.query('users', where: 'id = ?', whereArgs: [id]);
    if (rows.isEmpty) return null;
    return User.fromJson(rows.first);
  }

  /// Returns all users, ordered by `created_at` ascending.
  static Future<List<User>> getAllUsers() async {
    final db = await database;
    final rows = await db.query('users', orderBy: 'created_at');
    return rows.map((r) => User.fromJson(r)).toList();
  }

  /// Updates only the provided fields plus `updated_at`, then returns the updated user.
  static Future<User> updateUser(int id, Map<String, dynamic> updates) async {
    final db = await database;
    updates['updated_at'] = DateTime.now().toIso8601String();
    await db.update('users', updates, where: 'id = ?', whereArgs: [id]);
    // Must exist because update succeeded
    return (await getUser(id))!;
  }

  /// Deletes the user with the given `id`.
  static Future<void> deleteUser(int id) async {
    final db = await database;
    await db.delete('users', where: 'id = ?', whereArgs: [id]);
  }

  /// Returns the total number of users.
  static Future<int> getUserCount() async {
    final db = await database;
    final result = await db.rawQuery('SELECT COUNT(*) AS count FROM users');
    return firstIntValue(result) ?? 0;
  }

  /// Searches users by name or email (case-insensitive LIKE).
  static Future<List<User>> searchUsers(String query) async {
    final db = await database;
    final term = '%$query%';
    final rows = await db.query(
      'users',
      where: 'name LIKE ? OR email LIKE ?',
      whereArgs: [term, term],
    );
    return rows.map((r) => User.fromJson(r)).toList();
  }

  /// Deletes all rows in both tables.
  static Future<void> clearAllData() async {
    final db = await database;
    await db.delete('posts');
    await db.delete('users');
  }

  /// Closes the database connection.
  static Future<void> closeDatabase() async {
    if (_database != null) {
      await _database!.close();
      _database = null;
    }
  }

  /// Returns the full filesystem path to the database file.
  static Future<String> getDatabasePath() async {
    return join(await getDatabasesPath(), _dbName);
  }
}
