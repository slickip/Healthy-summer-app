import 'package:flutter/material.dart';
import '../../services/api_service.dart';
import 'package:intl/intl.dart';

class EditMealScreen extends StatefulWidget {
  final Map<String, dynamic> meal;
  const EditMealScreen({Key? key, required this.meal}) : super(key: key);

  @override
  State<EditMealScreen> createState() => _EditMealScreenState();
}

class _EditMealScreenState extends State<EditMealScreen> {
  final ApiService _apiService = ApiService();
  final _formKey = GlobalKey<FormState>();
  late TextEditingController _descController;
  late TextEditingController _calController;
  final _timeController = TextEditingController();
  DateTime? _mealTime;

  @override
  void initState() {
    super.initState();
    _descController = TextEditingController(
      text: widget.meal['description'] ?? '',
    );
    _calController = TextEditingController(
      text: widget.meal['calories'].toString(),
    );
    final raw = widget.meal['meal_time'];
    _mealTime = DateTime.tryParse(raw);
    _timeController.text = DateFormat(
      'yyyy-MM-ddTHH:mm:ss',
    ).format(_mealTime ?? DateTime.now());
  }

  Future<void> _selectDateTime() async {
    final date = await showDatePicker(
      context: context,
      initialDate: _mealTime ?? DateTime.now(),
      firstDate: DateTime(2023),
      lastDate: DateTime(2100),
    );
    if (date == null) return;

    final time = await showTimePicker(
      context: context,
      initialTime: TimeOfDay.fromDateTime(_mealTime ?? DateTime.now()),
    );
    if (time == null) return;

    setState(() {
      _mealTime = DateTime(
        date.year,
        date.month,
        date.day,
        time.hour,
        time.minute,
      );
      _timeController.text = DateFormat(
        'yyyy-MM-ddTHH:mm:ss',
      ).format(_mealTime!);
    });
  }

  Future<void> _saveChanges() async {
    if (_formKey.currentState?.validate() ?? false && _mealTime != null) {
      await _apiService.updateMeal(
        id: widget.meal['id'],
        description: _descController.text.trim(),
        calories: int.tryParse(_calController.text.trim()),
        mealTime: _mealTime!.toUtc().toIso8601String(),
      );
      Navigator.pop(context);
    }
  }

  InputDecoration _dec(String label) => InputDecoration(
    labelText: label,
    filled: true,
    fillColor: Colors.orange[50],
    border: OutlineInputBorder(borderRadius: BorderRadius.circular(10)),
  );

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.orange[50],
      appBar: AppBar(
        backgroundColor: Colors.orange[700],
        title: const Text('Edit Meal'),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Form(
          key: _formKey,
          child: Column(
            children: [
              TextFormField(
                controller: _descController,
                decoration: _dec('Description'),
                validator: (v) =>
                    v == null || v.isEmpty ? 'Enter description' : null,
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _calController,
                decoration: _dec('Calories'),
                keyboardType: TextInputType.number,
                validator: (v) =>
                    int.tryParse(v ?? '') == null ? 'Enter calories' : null,
              ),
              const SizedBox(height: 12),
              Row(
                children: [
                  Expanded(
                    child: TextFormField(
                      controller: _timeController,
                      readOnly: true,
                      decoration: _dec('Meal Time'),
                      validator: (v) =>
                          v == null || v.isEmpty ? 'Select time' : null,
                    ),
                  ),
                  const SizedBox(width: 8),
                  ElevatedButton.icon(
                    onPressed: _selectDateTime,
                    icon: const Icon(Icons.calendar_today),
                    label: const Text('Select date & time'),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: Colors.orange[600],
                      padding: const EdgeInsets.symmetric(
                        horizontal: 12,
                        vertical: 16,
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 24),
              ElevatedButton.icon(
                onPressed: _saveChanges,
                icon: const Icon(Icons.save),
                label: const Text('Save Changes'),
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.orange[700],
                  padding: const EdgeInsets.symmetric(
                    horizontal: 32,
                    vertical: 14,
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
