import 'package:flutter/material.dart';
import '../../services/api_service.dart';

class AddFoodScreen extends StatefulWidget {
  const AddFoodScreen({Key? key}) : super(key: key);

  @override
  State<AddFoodScreen> createState() => _AddFoodScreenState();
}

class _AddFoodScreenState extends State<AddFoodScreen> {
  final ApiService _apiService = ApiService();
  final _formKey = GlobalKey<FormState>();
  final _name = TextEditingController();
  final _cal = TextEditingController();
  final _p = TextEditingController();
  final _f = TextEditingController();
  final _c = TextEditingController();

  Future<void> _submit() async {
    if (_formKey.currentState?.validate() ?? false) {
      final name = _name.text.trim();
      final calories = double.tryParse(_cal.text.trim()) ?? 0;
      final proteins = double.tryParse(_p.text.trim()) ?? 0;
      final fats = double.tryParse(_f.text.trim()) ?? 0;
      final carbs = double.tryParse(_c.text.trim()) ?? 0;

      print(
        'Sending food data: Name="$name", Calories=$calories, Proteins=$proteins, Fats=$fats, Carbs=$carbs',
      );

      final success = await _apiService.createFood(
        name: name,
        caloriesPer100g: calories,
        proteins: proteins,
        fats: fats,
        carbs: carbs,
      );

      print('Create food result: $success');

      if (success) Navigator.pop(context);
    }
  }

  InputDecoration _dec(String label) => InputDecoration(
    labelText: label,
    filled: true,
    fillColor: Colors.orange[50],
    border: OutlineInputBorder(borderRadius: BorderRadius.circular(10)),
  );

  String? _validateNum(String? v) =>
      double.tryParse(v ?? '') == null ? 'Enter valid number' : null;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.orange[50],
      appBar: AppBar(
        backgroundColor: Colors.orange[700],
        title: const Text('Add Food'),
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
                onPressed: _submit,
                icon: const Icon(Icons.save),
                label: const Text('Save Food'),
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
