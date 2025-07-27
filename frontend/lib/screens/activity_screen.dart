import 'package:flutter/material.dart';
import '../services/api_service.dart';

class AddActivityScreen extends StatefulWidget {
  const AddActivityScreen({Key? key}) : super(key: key);

  @override
  _AddActivityScreenState createState() => _AddActivityScreenState();
}

class _AddActivityScreenState extends State<AddActivityScreen> {
  final _durationController = TextEditingController();
  String _selectedType = 'running';
  String _selectedIntensity = 'medium';
  bool _loading = false;
  String? _error;

  final api = ApiService();

  final List<String> _types = ['running', 'cycling', 'swimming', 'yoga'];
  final List<String> _intensities = ['low', 'medium', 'high'];

  final Map<String, int> _typeToId = {
    'running': 1,
    'cycling': 2,
    'swimming': 3,
    'yoga': 4,
  };

  Future<void> _submit() async {
    setState(() {
      _loading = true;
      _error = null;
    });

    final success = await api.createActivity(
      activityTypeId: _typeToId[_selectedType]!,
      duration: int.tryParse(_durationController.text) ?? 0,
      intensity: _selectedIntensity,
    );

    if (success) {
      if (!mounted) return;
      Navigator.pop(context, true);
    } else {
      setState(() {
        _error = 'Failed to create activity.';
      });
    }

    setState(() {
      _loading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.orange[700],
        title: const Text('Add Activity'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            DropdownButtonFormField<String>(
              value: _selectedType,
              decoration: const InputDecoration(labelText: 'Type'),
              items: _types
                  .map(
                    (type) => DropdownMenuItem(value: type, child: Text(type)),
                  )
                  .toList(),
              onChanged: (value) => setState(() {
                _selectedType = value!;
              }),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _durationController,
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(
                labelText: 'Duration (minutes)',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 12),
            DropdownButtonFormField<String>(
              value: _selectedIntensity,
              decoration: const InputDecoration(labelText: 'Intensity'),
              items: _intensities
                  .map((i) => DropdownMenuItem(value: i, child: Text(i)))
                  .toList(),
              onChanged: (value) => setState(() {
                _selectedIntensity = value!;
              }),
            ),
            const SizedBox(height: 20),
            if (_error != null)
              Text(_error!, style: const TextStyle(color: Colors.red)),
            ElevatedButton(
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.orange[700],
                padding: const EdgeInsets.symmetric(vertical: 16),
              ),
              onPressed: _loading ? null : _submit,
              child: _loading
                  ? const CircularProgressIndicator(color: Colors.white)
                  : const Text('Save'),
            ),
          ],
        ),
      ),
    );
  }
}
