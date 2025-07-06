import 'package:flutter/material.dart';
import '../services/api_service.dart';

class ActivitiesScreen extends StatefulWidget {
  const ActivitiesScreen({Key? key}) : super(key: key);

  @override
  _ActivitiesScreenState createState() => _ActivitiesScreenState();
}

class _ActivitiesScreenState extends State<ActivitiesScreen> {
  final api = ApiService();
  bool _loading = true;
  List<dynamic> _activities = [];

  @override
  void initState() {
    super.initState();
    _loadActivities();
  }

  Future<void> _loadActivities() async {
    final data = await api.getActivities();
    setState(() {
      _activities = data ?? [];
      _loading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.orange[700],
        title: const Text('Your Activities'),
      ),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : ListView.builder(
              itemCount: _activities.length,
              itemBuilder: (context, index) {
                final a = _activities[index];
                return ListTile(
                  title: Text('${a['type']} (${a['duration']} min)'),
                  subtitle: Text(
                    'Calories: ${a['calories']} — ${a['intensity']}',
                  ),
                );
              },
            ),
      floatingActionButton: FloatingActionButton(
        backgroundColor: Colors.orange[700],
        child: const Icon(Icons.add),
        onPressed: () async {
          final result = await Navigator.pushNamed(context, '/add_activity');
          if (result == true) _loadActivities();
        },
      ),
    );
  }
}
