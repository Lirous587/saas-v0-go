INSERT INTO public.plans (id, name, price, description) VALUES
  (1, 'Free',        0.00, '免费计划'),
  (2, 'Caring',      0.00, '爱心计划'),
  (3, 'Professional',5.00, '专业计划')
ON CONFLICT (name) DO NOTHING;