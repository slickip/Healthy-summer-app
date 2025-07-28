import 'package:flutter/material.dart';
import '../../services/api_service.dart';
import 'package:intl/intl.dart';

class WaterScreen extends StatefulWidget {
  const WaterScreen({Key? key}) : super(key: key);

  @override
  State<WaterScreen> createState() => _WaterScreenState();
}

class _WaterScreenState extends State<WaterScreen> {
  final ApiService _apiService = ApiService();
  List<dynamic> _logs = [];

  @override
  void initState() {
    super.initState();
    _loadLogs();
  }

  Future<void> _loadLogs() async {
    final logs = await _apiService.getWaterLogs();
    if (logs != null) {
      setState(() {
        _logs = logs;
      });
    }
  }

  Future<void> _deleteLog(int id) async {
    await _apiService.deleteWaterLog(id);
    _loadLogs();
  }

  String _formatDate(String iso) {
    try {
      return DateFormat('yyyy-MM-dd HH:mm').format(DateTime.parse(iso));
    } catch (_) {
      return iso;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.orange[50],
      appBar: AppBar(
        title: const Text('Water Log'),
        backgroundColor: Colors.orange[700],
      ),
      body: ListView.builder(
        padding: const EdgeInsets.all(8),
        itemCount: _logs.length,
        itemBuilder: (context, index) {
          final log = _logs[index];
          return Card(
            color: Colors.white,
            elevation: 2,
            margin: const EdgeInsets.symmetric(vertical: 6, horizontal: 8),
            child: ListTile(
              title: Text('${log['volume_ml']} ml'),
              subtitle: Text(_formatDate(log['logged_at'] ?? '')),
              trailing: IconButton(
                icon: const Icon(Icons.delete),
                onPressed: () => _deleteLog(log['id']),
              ),
            ),
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () =>
            Navigator.pushNamed(context, '/add_water').then((_) => _loadLogs()),
        backgroundColor: Colors.orange[700],
        child: const Icon(Icons.add),
      ),
    );
  }
}
