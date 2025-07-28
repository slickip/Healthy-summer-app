import 'package:flutter/material.dart';
import '../../services/api_service.dart';

class EditFoodScreen extends StatefulWidget {
  final Map<String, dynamic> food;
  const EditFoodScreen({Key? key, required this.food}) : super(key: key);

  @override
  State<EditFoodScreen> createState() => _EditFoodScreenState();
}

class _EditFoodScreenState extends State<EditFoodScreen> {
  final ApiService _apiService = ApiService();
  final _formKey = GlobalKey<FormState>();
  late TextEditingController _name;
  late TextEditingController _cal;
  late TextEditingController _p;
  late TextEditingController _f;
  late TextEditingController _c;

  @override
  void initState() {
    super.initState();
    _name = TextEditingController(text: widget.food['name'] ?? '');
    _cal = TextEditingController(
      text: widget.food['callories_per_100g']?.toString() ?? '',
    );
    _p = TextEditingController(text: widget.food['proteins']?.toString() ?? '');
    _f = TextEditingController(text: widget.food['fats']?.toString() ?? '');
    _c = TextEditingController(text: widget.food['carbs']?.toString() ?? '');
  }

  @override
  void dispose() {
    _name.dispose();
    _cal.dispose();
    _p.dispose();
    _f.dispose();
    _c.dispose();
    super.dispose();
  }

  InputDecoration _dec(String label) => InputDecoration(
    labelText: label,
    filled: true,
    fillColor: Colors.orange[50],
    border: OutlineInputBorder(borderRadius: BorderRadius.circular(10)),
  );

  String? _validateNum(String? v) =>
      double.tryParse(v ?? '') == null ? 'Enter valid number' : null;

  Future<void> _saveChanges() async {
    if (_formKey.currentState?.validate() ?? false) {
      await _apiService.updateFood(
        id: widget.food['id'],
        name: _name.text.trim(),
        caloriesPer100g: double.tryParse(_cal.text.trim()) ?? 0,
        proteins: double.tryParse(_p.text.trim()) ?? 0,
        fats: double.tryParse(_f.text.trim()) ?? 0,
        carbs: double.tryParse(_c.text.trim()) ?? 0,
      );
      Navigator.pop(context);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.orange[50],
      appBar: AppBar(
        backgroundColor: Colors.orange[700],
        title: const Text('Edit Food'),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Form(
          key: _formKey,
          child: Column(
            children: [
              TextFormField(
                controller: _name,
                decoration: _dec('Name'),
                validator: (v) => v == null || v.isEmpty ? 'Enter name' : null,
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _cal,
                decoration: _dec('Calories per 100g'),
                keyboardType: TextInputType.number,
                validator: _validateNum,
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _p,
                decoration: _dec('Proteins (g)'),
                keyboardType: TextInputType.number,
                validator: _validateNum,
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _f,
                decoration: _dec('Fats (g)'),
                keyboardType: TextInputType.number,
                validator: _validateNum,
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _c,
                decoration: _dec('Carbohydrates (g)'),
                keyboardType: TextInputType.number,
                validator: _validateNum,
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
