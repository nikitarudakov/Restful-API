package voteDao

import "git.foxminded.ua/foxstudent106092/user-management/internal/business/model"

func MarshalTogRPCRating(r *model.Rating) *Rating {
	return &Rating{
		Target: r.Target,
		Rating: r.Rating,
	}
}

func UnmarshalToRating(r *Rating) *model.Rating {
	return &model.Rating{
		Target: r.Target,
		Rating: r.Rating,
	}
}

func MarshalTogRPCVote(v *model.Vote) *Vote {
	return &Vote{
		Sender:  v.Sender,
		Target:  v.Target,
		Vote:    v.Vote,
		VotedAt: v.VotedAt,
	}
}

func MarshalTogRPCInsertResult(ir *model.InsertResult) *InsertResult {
	var insertResult InsertResult

	insertResult.ObjectID = []byte(ir.Id.Hex())

	return &insertResult
}

func MarshalTogRPCVoteFilter(v *model.Vote, isTarget bool, isSender bool) *VoteFilter {
	return &VoteFilter{
		Vote: &Vote{
			Sender: v.Sender,
			Target: v.Target,
		},
		IsTarget: isTarget,
		IsSender: isSender,
	}
}
