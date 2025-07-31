import 'package:flutter/material.dart';
import '../../services/api_service.dart';

class CreateChallengeScreen extends StatefulWidget {
  const CreateChallengeScreen({super.key});

  @override
  State<CreateChallengeScreen> createState() => _CreateChallengeScreenState();
}

class _CreateChallengeScreenState extends State<CreateChallengeScreen> {
  final ApiService api = ApiService();
  final TextEditingController _titleController = TextEditingController();
  final TextEditingController _descriptionController = TextEditingController();
  final TextEditingController _goalValueController = TextEditingController();

  Future<void> submit() async {
    final title = _titleController.text;
    final desc = _descriptionController.text;
    final goal = int.tryParse(_goalValueController.text);

    if (title.isEmpty || desc.isEmpty || goal == null) {
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(const SnackBar(content: Text('Fill in all fields')));
      return;
    }

    final success = await api.createChallenge(
      title: title,
      description: desc,
      goalValue: goal,
    );

    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(success ? 'Challenge created!' : 'Error creating'),
      ),
    );

    if (success) {
      Navigator.pop(context);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.orange[50],
      appBar: AppBar(
        backgroundColor: Colors.orange[700],
        title: const Text(
          'Create Challenge',
          style: TextStyle(color: Colors.white),
        ),
      ),
      body: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          children: [
            TextField(
              controller: _titleController,
              decoration: const InputDecoration(
                labelText: 'Title',
                filled: true,
                fillColor: Colors.white,
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _descriptionController,
              decoration: const InputDecoration(
                labelText: 'Description',
                filled: true,
                fillColor: Colors.white,
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _goalValueController,
              decoration: const InputDecoration(
                labelText: 'Goal Value',
                filled: true,
                fillColor: Colors.white,
                border: OutlineInputBorder(),
              ),
              keyboardType: TextInputType.number,
            ),
            const SizedBox(height: 20),
            ElevatedButton.icon(
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.orange[700],
              ),
              onPressed: submit,
              icon: const Icon(Icons.check),
              label: const Text('Create'),
            ),
          ],
        ),
      ),
    );
  }
}
