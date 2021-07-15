package mysql

// func (db *DB) Approve(ctx context.Context, p model.Package) error {
// 	if err := p.Validate(); err != nil {
// 		return err
// 	}

// 	_, err := db.Exec(
// 		`UPDATE packets set approved=!approved WHERE packname=$1 AND majver=$2 AND minver=$3 AND arch=$4`,
// 		p.Name,
// 		p.VersionMajor,
// 		p.VersionMinor,
// 		p.Architecture,
// 	)
// 	return err
// }

// func (db *DB) Remove(ctx context.Context, p model.Package) error {
// 	if err := p.Validate(); err != nil {
// 		return err
// 	}

// 	_, err := db.Exec(
// 		`DELETE FROM packets WHERE packname=$1 AND majver=$2 AND minver=$3 AND arch=$4 LIMIT 1`,
// 		p.Name,
// 		p.VersionMajor,
// 		p.VersionMinor,
// 		p.Architecture,
// 	)
// 	return err
// }
