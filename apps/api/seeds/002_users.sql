-- Password: password123
INSERT INTO users (id, name, email, password_hash) VALUES
(
  'b1000000-0000-4000-8000-000000000001',
  'Aman',
  'aman@ylx.dev',
  '$argon2id$v=19$m=65536,t=3,p=2$ZlpIeWUqwDTEqrSO$tjWYwf3f3NkvQYb4ZjTSAXmm9sbpkrtMiNp4/XGN03Q'
),
(
  'b1000000-0000-4000-8000-000000000002',
  'Sujoy',
  'sujoy@ylx.dev',
  '$argon2id$v=19$m=65536,t=3,p=2$SNFQBzn4MzLTI7Tk$K6HGeXafxrhaA28vFO0DeUbeUpAx+Ph4STLUQMQ9kK4'
),
(
  'b1000000-0000-4000-8000-000000000003',
  'Akash',
  'akash@ylx.dev',
  '$argon2id$v=19$m=65536,t=3,p=2$zGmTbN96Ds+uXYBs$FrOH74xTC1Ax3aQlhZfLBtrqcCshndbNjCBGMNk4xaM'
)
ON CONFLICT (id) DO NOTHING;
